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

// StoreNonce stores nonce per chainID
func (ns *NonceStore) StoreNonce(chainID *big.Int, nonce *big.Int) error {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:nonce", chainID.String())
	key.WriteString(keyS)

	err := ns.db.SetByKey(key.Bytes(), nonce.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// GetLastNonce queries the blockstore and returns latest nonce
func (ns *NonceStore) GetLastNonce(chainID *big.Int, nonce *big.Int) (*big.Int, error) {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:nonce", chainID.String())
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
