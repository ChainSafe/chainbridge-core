package transactor

import (
	"context"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/imdario/mergo"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type TransactOptions struct {
	GasLimit uint64
	GasPrice *big.Int
	Value    *big.Int
}

var DefaultTransactionOptions = TransactOptions{
	GasLimit: 2000000,
	GasPrice: big.NewInt(0),
	Value:    big.NewInt(0),
}

func MergeTransactionOptions(primary TransactOptions, additional TransactOptions) TransactOptions {
	if err := mergo.Merge(&primary, additional); err != nil {
		log.Fatal().Msg("Unable to merge")
		return TransactOptions{}
	}
	return primary
}

type Transactor interface {
	Transact(to *common.Address, data []byte, opts TransactOptions) (*common.Hash, error)
}

type signAndSendTransactor struct {
	TxFabric       calls.TxFabric
	gasPriceClient calls.GasPricer
	client         calls.ClientDispatcher
}

func NewSignAndSendTransactor(txFabric calls.TxFabric, gasPriceClient calls.GasPricer, client calls.ClientDispatcher) Transactor {
	return &signAndSendTransactor{
		TxFabric:       txFabric,
		gasPriceClient: gasPriceClient,
		client:         client,
	}
}

func (t *signAndSendTransactor) Transact(to *common.Address, data []byte, opts TransactOptions) (*common.Hash, error) {
	defer t.client.UnlockNonce()
	t.client.LockNonce()
	n, err := t.client.UnsafeNonce()
	if err != nil {
		return &common.Hash{}, nil
	}
	opts = MergeTransactionOptions(opts, DefaultTransactionOptions)
	gp := []*big.Int{opts.GasPrice}
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice()
		if err != nil {
			return &common.Hash{}, err
		}
	}

	tx, err := t.TxFabric(n.Uint64(), to, opts.Value, opts.GasLimit, gp, data)
	if err != nil {
		return &common.Hash{}, err
	}
	h, err := t.client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		log.Error().Err(err).Msg("SIGN_AND_SEND")
		return &common.Hash{}, err
	}
	_, err = t.client.WaitAndReturnTxReceipt(h)
	if err != nil {
		return &common.Hash{}, err
	}
	err = t.client.UnsafeIncreaseNonce()
	if err != nil {
		return &common.Hash{}, err
	}
	return &h, nil
}
