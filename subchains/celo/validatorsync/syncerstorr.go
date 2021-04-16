//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"math/big"
	"time"

	"github.com/celo-org/celo-blockchain/consensus/istanbul"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	latestKnownBlockKey = "latestKnownBlock"
)

func NewValidatorsStore(db *leveldb.DB) *ValidatorsStore {
	return &ValidatorsStore{db: db}
}

type ValidatorsStore struct {
	db *leveldb.DB
}

// GetLatestKnownBlock returns block number of latest parsed EpochLastBlock for provided chainID. If DB is empty returns 0.
// Should always be last block in epoch.
func (db *ValidatorsStore) GetLatestKnownEpochLastBlock(chainID uint8) (*big.Int, error) {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return nil, err
	}
	key.WriteString(latestKnownBlockKey)
	data, err := db.db.Get(key.Bytes(), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return big.NewInt(0), nil
		}
		return nil, err
	}
	v := big.NewInt(0)
	v.SetBytes(data)
	return v, nil
}

// Atomically sets block and validators as related KV to underlying DB backend
func (db *ValidatorsStore) SetValidatorsForBlock(block *big.Int, validators []*istanbul.ValidatorData, chainID uint8) error {
	byteValidators := &bytes.Buffer{}
	enc := gob.NewEncoder(byteValidators)
	err := enc.Encode(validators)
	if err != nil {
		return err
	}
	tx, err := db.db.OpenTransaction()
	if err != nil {
		return err
	}
	key := new(bytes.Buffer)
	err = binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return err
	}
	key.Write(block.Bytes())
	err = tx.Put(key.Bytes(), byteValidators.Bytes(), nil)
	if err != nil {
		tx.Discard()
		return err
	}
	err = db.setLatestKnownEpochLastBlockWithTransaction(block, chainID, tx)
	if err != nil {
		tx.Discard()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Discard()
		return err
	}
	return nil
}

func (db *ValidatorsStore) GetValidatorsForBlock(block *big.Int, chainID uint8) ([]*istanbul.ValidatorData, error) {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return nil, err
	}
	key.Write(block.Bytes())
	res, err := db.db.Get(key.Bytes(), nil)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	b.Write(res)
	dec := gob.NewDecoder(b)
	dataArr := make([]*istanbul.ValidatorData, 0)
	err = dec.Decode(&dataArr)
	if err != nil {
		return nil, err
	}
	return dataArr, nil
}

func (db *ValidatorsStore) setLatestKnownEpochLastBlockWithTransaction(block *big.Int, chainID uint8, transaction *leveldb.Transaction) error {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return err
	}
	key.WriteString(latestKnownBlockKey)
	err = transaction.Put(key.Bytes(), block.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

var ErrNoBlockInStore = errors.New("no corresponding validators for provided block number")

func (db *ValidatorsStore) GetAPKForBlock(block *big.Int, chainID uint8, epochSize uint64) ([]byte, error) {
	for i := 0; i <= 10; i++ {
		vals, err := db.GetValidatorsForBlock(computeLastBlockOfEpochForProvidedBlock(block, epochSize), chainID)
		if err != nil {
			if errors.Is(err, leveldb.ErrNotFound) {
				time.Sleep(5 * time.Second)
				continue
			}
			return nil, err
		}
		pk, err := aggregatePublicKeys(vals)
		if err != nil {
			return nil, err
		}
		return pk.Serialize()
	}
	return nil, ErrNoBlockInStore
}

// Closes connection to underlying DB backend
func (db *ValidatorsStore) Close() error {
	if err := db.db.Close(); err != nil {
		return err
	}
	return nil
}
