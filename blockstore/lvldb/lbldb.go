package lvldb

import (
	"math/big"

	"github.com/syndtr/goleveldb/leveldb"
)

type LVLDBBlockStore struct {
	db *leveldb.DB
}

func NewLVLDBBlockStore(db *leveldb.DB) *LVLDBBlockStore {
	return &LVLDBBlockStore{db: db}
}

func (db *LVLDBBlockStore) StoreBlock(block *big.Int, chainID uint8) error {

}

func (db *LVLDBBlockStore) GetLastStoredBlock(chainID uint8) error {

}

// Closes connection to underlying DB backend
func (db *LVLDBBlockStore) Close() error {
	if err := db.db.Close(); err != nil {
		return err
	}
	return nil
}
