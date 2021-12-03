package store

import "errors"

var (
	ErrNotFound = errors.New("key not found")
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
