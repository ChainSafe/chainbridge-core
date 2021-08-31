package evmtypes

import (
	"math/big"
)

// GasPricer abstract model for providers of gasPrice. Depending on the implementation array of gasprices could be differ.
//In case of static implementation (pre EIP-1559) only one elemnt will be returned. For EIP1559 implementations 3 elements will be returned
type GasPricer interface {
	GasPrice() ([]*big.Int, error)
}
