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

type NonceStore struct {
	db KeyValueReaderWriter
}

func NewNonceStore(db KeyValueReaderWriter) *NonceStore {
	return &NonceStore{
		db: db,
	}
}

// StoreNonce stores nonce per chainID
func (ns *NonceStore) StoreNonce(chainID *big.Int, nonce *big.Int) error {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%d:nonce", chainID.Int64())
	key.WriteString(keyS)

	err := ns.db.SetByKey(key.Bytes(), nonce.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// GetNonce queries the blockstore and returns latest nonce
func (ns *NonceStore) GetNonce(chainID *big.Int) (*big.Int, error) {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%d:nonce", chainID.Int64())
	key.WriteString(keyS)

	v, err := ns.db.GetByKey(key.Bytes())
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return big.NewInt(0), nil
		}
		return nil, err
	}

	block := big.NewInt(0).SetBytes(v)
	return block, nil
}
