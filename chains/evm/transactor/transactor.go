package transactor

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransactOptions struct {
	GasLimit *big.Int
	GasPrice *big.Int
	Value    *big.Int
	Nonce    *big.Int
	ChainID  uint8
	Priority Priority
}

type Priority string

const (
	LowPriority  Priority = "low"
	MedPriority           = "medium"
	HighPriority          = "high"
)

type Transactor interface {
	Transact(to *common.Address, data []byte, opts TransactOptions) (common.Hash, error)
}
