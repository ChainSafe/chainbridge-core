package transaction

import (
	"context"
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxFabric func(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrices []*big.Int, data []byte) (client.CommonTransaction, error)

type TX struct {
	tx *types.Transaction
}

// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
// but return raw byte representation of transaction to be compatible and interchangeable between different go-ethereum forks
// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (a *TX) RawWithSignature(signer client.Signer, domainID *big.Int) ([]byte, error) {
	opts, err := newTransactorWithChainID(signer, domainID)
	if err != nil {
		return nil, err
	}
	tx, err := opts.Signer(signer.CommonAddress(), a.tx)
	if err != nil {
		return nil, err
	}
	a.tx = tx

	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// NewTransaction is the ethereum transaction constructor
func NewTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrices []*big.Int, data []byte) (client.CommonTransaction, error) {
	// If there is more than one gas price returned we are sending with DynamicFeeTx's
	if len(gasPrices) > 1 {
		return newDynamicFeeTransaction(nonce, to, amount, gasLimit, gasPrices[0], gasPrices[1], data), nil
	} else {
		return newTransaction(nonce, to, amount, gasLimit, gasPrices[0], data), nil
	}
}

func newDynamicFeeTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasTipCap *big.Int, gasFeeCap *big.Int, data []byte) *TX {
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        to,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		Value:     amount,
		Data:      data,
	})
	return &TX{tx: tx}
}

func newTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *TX {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
	return &TX{tx: tx}
}

func (a *TX) Hash() common.Hash {
	return a.tx.Hash()
}

// newTransactorWithChainID is a utility method to easily create a transaction signer
// for an client.Signer.
// Mostly copies bind.NewKeyedTransactorWithChainID but sings with the provided signer
// instead of a privateKey
func newTransactorWithChainID(s client.Signer, chainID *big.Int) (*bind.TransactOpts, error) {
	keyAddr := s.CommonAddress()
	if chainID == nil {
		return nil, bind.ErrNoChainID
	}
	signer := types.LatestSignerForChainID(chainID)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, bind.ErrNotAuthorized
			}
			signature, err := s.Sign(signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
		Context: context.Background(),
	}, nil
}
