package transactor

import (
	"math/big"

	"github.com/imdario/mergo"

	"github.com/ethereum/go-ethereum/common"
)

type TransactOptions struct {
	GasLimit uint64
	GasPrice *big.Int
	Value    *big.Int
	Nonce    *big.Int
	ChainID  *big.Int
	Priority string
}

func MergeTransactionOptions(primary *TransactOptions, additional *TransactOptions) error {
	if err := mergo.Merge(primary, additional); err != nil {
		return err
	}

	return nil
}

type Transactor interface {
	Transact(to *common.Address, data []byte, opts TransactOptions) (*common.Hash, error)
}
