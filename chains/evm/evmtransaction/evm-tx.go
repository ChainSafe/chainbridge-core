package evmtransaction

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type TX struct {
	types.Transaction
}

// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
// but return raw byte representation of transaction to be compatible and interchangeable between different go-ethereum forks
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
