package evmtransaction

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
)

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

func (txOpts *TransactOpts) SetNonce(nonce *big.Int) {
	txOpts.opts.Nonce = nonce
	log.Debug().Msgf("Nonce inside SetNonce: %v", txOpts.opts.Nonce)
}

func (txOpts *TransactOpts) SetGasPrices(gasPrice *big.Int, gasTipCap *big.Int, gasFeeCap *big.Int) {
	txOpts.opts.GasPrice = gasPrice
	txOpts.opts.GasTipCap = gasTipCap
	txOpts.opts.GasFeeCap = gasFeeCap
}

func (txOpts *TransactOpts) SetGasLimit(gasLimit uint64) {
	txOpts.opts.GasLimit = gasLimit
}

func (txOpts *TransactOpts) Signer() bind.SignerFn {
	return txOpts.opts.Signer
}

func (txOpts *TransactOpts) GasPrice() *big.Int {
	return txOpts.opts.GasPrice
}

func (txOpts *TransactOpts) GasTipCap() *big.Int {
	return txOpts.opts.GasTipCap
}

func (txOpts *TransactOpts) GasFeeCap() *big.Int {
	return txOpts.opts.GasFeeCap
}

func (txOpts *TransactOpts) GasLimit() uint64 {
	return txOpts.opts.GasLimit
}

func (txOpts *TransactOpts) Nonce() *big.Int {
	return txOpts.opts.Nonce
}
