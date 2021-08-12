package evmtransaction

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
)

// SignerFn is a signer function callback when a contract requires a method to
// sign the transaction before submission.
//type SignerFn func(common.Address, *types.Transaction) (*types.Transaction, error)

type CommonTransactOpts interface {
	SetNonce(*big.Int)
	SetGasPrices(*big.Int, *big.Int, *big.Int)
	SetGasLimit(uint64)

	Nonce() *big.Int
	GasPrice() *big.Int
	GasTipCap() *big.Int
	GasFeeCap() *big.Int
	GasLimit() uint64
}

type EVMTransactor interface {
	CommonTransactOpts
	// NOTE: should we declare this function as native Signer type rather than using bind package
	Signer() bind.SignerFn
}

type TransactOpts struct {
	opts *bind.TransactOpts
}

func NewOpts(key *ecdsa.PrivateKey, chainID *big.Int) (*TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}
	return &TransactOpts{opts: opts}, nil
}

func (t *TransactOpts) SetNonce(nonce *big.Int) {
	t.opts.Nonce = nonce
	log.Debug().Msgf("Nonce inside SetNonce: %v", t.opts.Nonce)
}

func (t *TransactOpts) SetGasPrices(gasPrice *big.Int, gasTipCap *big.Int, gasFeeCap *big.Int) {
	t.opts.GasPrice = gasPrice
	t.opts.GasTipCap = gasTipCap
	t.opts.GasFeeCap = gasFeeCap
}

func (t *TransactOpts) SetGasLimit(gasLimit uint64) {
	t.opts.GasLimit = gasLimit
}

func (t *TransactOpts) Signer() bind.SignerFn {
	return t.opts.Signer
}

func (t *TransactOpts) GasPrice() *big.Int {
	return t.opts.GasPrice
}

func (t *TransactOpts) GasTipCap() *big.Int {
	return t.opts.GasTipCap
}

func (t *TransactOpts) GasFeeCap() *big.Int {
	return t.opts.GasFeeCap
}

func (t *TransactOpts) GasLimit() uint64 {
	return t.opts.GasLimit
}

func (t *TransactOpts) Nonce() *big.Int {
	return t.opts.Nonce
}
