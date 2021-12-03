// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package store

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/syndtr/goleveldb/leveldb"
)

type BlockStore struct {
	db KeyValueReaderWriter
}

func NewBlockStore(db KeyValueReaderWriter) *BlockStore {
	return &BlockStore{
		db: db,
	}
}

// StoreBlock stores block number per domainID into blockstore
func (bs *BlockStore) StoreBlock(block *big.Int, domainID uint8) error {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%d:block", domainID)
	key.WriteString(keyS)

	err := bs.db.SetByKey(key.Bytes(), block.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// GetLastStoredBlock queries the blockstore and returns latest known block
func (bs *BlockStore) GetLastStoredBlock(domainID uint8) (*big.Int, error) {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%d:block", domainID)
	key.WriteString(keyS)

	v, err := bs.db.GetByKey(key.Bytes())
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return big.NewInt(0), nil
		}
		return nil, err
	}

	block := big.NewInt(0).SetBytes(v)
	return block, nil
}

// GetStartBlock queries the blockstore for the latest known block. If the latest block is
// greater than configured startBlock, then startBlock is replaced with the latest known block.
func (bs *BlockStore) GetStartBlock(domainID uint8, startBlock *big.Int, latest bool, fresh bool) (*big.Int, error) {
	if latest {
		return nil, nil
	}

	if fresh {
		return startBlock, nil
	}

	latestBlock, err := bs.GetLastStoredBlock(domainID)
	if err != nil {
		return nil, err
	}

	if latestBlock.Cmp(startBlock) == 1 {
		return latestBlock, nil
	} else {
		return startBlock, nil
	}
}
