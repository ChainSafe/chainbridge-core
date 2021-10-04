// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package lvldb

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type LVLDB struct {
	db *leveldb.DB
}

func NewLvlDB(path string) (*LVLDB, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "levelDB.OpenFile fail")
	}
	return &LVLDB{db: ldb}, nil
}

func (db *LVLDB) GetByKey(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *LVLDB) SetByKey(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *LVLDB) Close() error {
	return db.db.Close()
}
