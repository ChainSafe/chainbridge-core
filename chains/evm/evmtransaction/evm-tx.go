package evmtransaction

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type TX struct {
	tx *types.Transaction
}

// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
// but return raw byte representation of transaction to be compatible and interchangeable between different go-ethereum forks
// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (a *TX) RawWithSignature(key *ecdsa.PrivateKey, chainID *big.Int) ([]byte, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}
	tx, err := opts.Signer(crypto.PubkeyToAddress(key.PublicKey), a.tx)
	if err != nil {
		return nil, err
	}
	a.tx = tx
	rawTX, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return rawTX, nil
}

func NewTransaction(opts *bind.TransactOpts, to common.Address, amount *big.Int, data []byte) *TX {
	var tx *types.Transaction
	if opts.GasFeeCap != nil {
		tx = types.NewTx(&types.DynamicFeeTx{
			Nonce:     opts.Nonce.Uint64(),
			To:        &to,
			GasFeeCap: opts.GasFeeCap,
			GasTipCap: opts.GasTipCap,
			Gas:       opts.GasLimit,
			Value:     amount,
			Data:      data,
		})
	} else {
		tx = types.NewTransaction(opts.Nonce.Uint64(), to, amount, opts.GasLimit, opts.GasPrice, data)
	}
	return &TX{tx: tx}
}

func (a *TX) Hash() common.Hash {
	return a.tx.Hash()
}
