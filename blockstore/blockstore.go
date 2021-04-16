package blockstore

import (
	"bytes"
	"fmt"
	"math/big"
)

type KeyValueReaderWriter interface {
	GetByKey(key []byte) ([]byte, error)
	SetByKey(key []byte, value []byte) error
}

func StoreBlock(db KeyValueReaderWriter, block *big.Int, chainID uint8) error {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:block", string(chainID))
	key.WriteString(keyS)
	err := db.SetByKey(key.Bytes(), block.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func GetLastStoredBlock(db KeyValueReaderWriter, chainID uint8) (*big.Int, error) {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:block", string(chainID))
	key.WriteString(keyS)
	v, err := db.GetByKey(key.Bytes())
	if err != nil {
		return nil, err
	}
	block := big.NewInt(0).SetBytes(v)
	return block, nil
}
