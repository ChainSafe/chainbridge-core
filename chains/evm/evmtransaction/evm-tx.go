package evmtransaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type TX struct {
	types.Transaction
}

// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
// but return raw byte representation of transaction to be compatible and interchangeable between different go-ethereum forks
// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (a *TX) RawWithSignature(signer types.Signer, sig []byte) ([]byte, error) {
	tx, err := a.Transaction.WithSignature(signer, sig)
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *TX {
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
	return &TX{*tx}
}
