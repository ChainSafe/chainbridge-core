package blockstore

import (
	"bytes"
	"fmt"
	"math/big"
)

type KeyValueDB interface {
	GetByKey(key []byte) ([]byte, error)
	SetByKey(key []byte, value []byte) error
	Close() error
}

type BlockStore struct {
	backend KeyValueDB
}

func NewBlockStore(db KeyValueDB) (*BlockStore, error) {
	return &BlockStore{backend: db}, nil
}

func (db *BlockStore) StoreBlock(block *big.Int, chainID uint8) error {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:block", string(chainID))
	key.WriteString(keyS)
	err := db.backend.SetByKey(key.Bytes(), block.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (db *BlockStore) GetLastStoredBlock(chainID uint8) (*big.Int, error) {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:block", string(chainID))
	key.WriteString(keyS)
	v, err := db.backend.GetByKey(key.Bytes())
	if err != nil {
		return nil, err
	}
	block := big.NewInt(0).SetBytes(v)
	return block, nil
}

// Closes connection to underlying DB backend
func (db *BlockStore) Close() error {
	if err := db.backend.Close(); err != nil {
		return err
	}
	return nil
}
