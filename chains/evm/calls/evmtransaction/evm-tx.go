package evmtransaction

import (
	"crypto/ecdsa"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type TX struct {
	tx *types.Transaction
}

// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
// but return raw byte representation of transaction to be compatible and interchangeable between different go-ethereum forks
// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (a *TX) RawWithSignature(key *ecdsa.PrivateKey, domainID *big.Int) ([]byte, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(key, domainID)
	if err != nil {
		return nil, err
	}
	tx, err := opts.Signer(crypto.PubkeyToAddress(key.PublicKey), a.tx)
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
func NewTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrices []*big.Int, data []byte) (evmclient.CommonTransaction, error) {
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
