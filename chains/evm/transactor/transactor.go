package transactor

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/imdario/mergo"
)

type TransactOptions struct {
	GasLimit *big.Int
	GasPrice *big.Int
	Value    *big.Int
	Nonce    *big.Int
	ChainID  uint8
	Priority string
}

type Transactor interface {
	Transact(to *common.Address, data []byte, opts TransactOptions) (common.Hash, error)
}

func MergeTransactionOptions(primary *TransactOptions, additional *TransactOptions) error {
	if err := mergo.Merge(primary, additional); err != nil {
		return err
	}

	return nil
}
