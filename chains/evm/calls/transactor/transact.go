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
	Priority uint8
}

// to save on data, we encode uin8 for transaction priority
var TxPriorities = map[string]uint8{
	"none":   0,
	"slow":   1,
	"medium": 2,
	"fast":   3,
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
