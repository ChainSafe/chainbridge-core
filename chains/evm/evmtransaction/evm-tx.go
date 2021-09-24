package evmtransaction

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"

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
func (a *TX) RawWithSignature(key *ecdsa.PrivateKey, chainId *big.Int) ([]byte, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(key, chainId)
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

// NewTransaction is the
func NewTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPricer calls.GasPricer, data []byte) evmclient.CommonTransaction {
	gasPrices, err := gasPricer.GasPrice()
	if err != nil {
		return nil
	}
	// If there is more than one gas price returned we are sending with DynamicFeeTx's
	if gasPrices[1] == nil {
		return newTransaction(nonce, to, amount, gasLimit, gasPrices[0], data)
	} else {
		return newDynamicFeeTransaction(nonce, to, amount, gasLimit, gasPrices[0], gasPrices[1], data)
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
	var tx *types.Transaction
	if to == nil {
		tx = types.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)
	} else {
		tx = types.NewTransaction(nonce, *to, amount, gasLimit, gasPrice, data)
	}
	return &TX{tx: tx}
}

func (a *TX) Hash() common.Hash {
	return a.tx.Hash()
}
