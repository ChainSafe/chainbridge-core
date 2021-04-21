package blockstore

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/syndtr/goleveldb/leveldb"
)

type KeyValueReaderWriter interface {
	KeyValueReader
	KeyValueWriter
}

type KeyValueReader interface {
	GetByKey(key []byte) ([]byte, error)
}

type KeyValueWriter interface {
	SetByKey(key []byte, value []byte) error
}

var (
	ErrNotFound = errors.New("key not found")
)

func StoreBlock(db KeyValueWriter, block *big.Int, chainID uint8) error {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:block", string(chainID))
	key.WriteString(keyS)
	err := db.SetByKey(key.Bytes(), block.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func GetLastStoredBlock(db KeyValueReader, chainID uint8) (*big.Int, error) {
	key := bytes.Buffer{}
	keyS := fmt.Sprintf("chain:%s:block", string(chainID))
	key.WriteString(keyS)
	v, err := db.GetByKey(key.Bytes())
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return big.NewInt(0), nil
		}
		return nil, err
	}
	block := big.NewInt(0).SetBytes(v)
	return block, nil
}
