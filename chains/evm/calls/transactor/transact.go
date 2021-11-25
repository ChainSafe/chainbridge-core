package transactor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type TransactOptions struct {
	GasLimit uint64
	GasPrice *big.Int
	Value    *big.Int
}

func NewDefaultTransactOptions() TransactOptions {
	return TransactOptions{
		GasLimit: 2000000,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
	}
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

	gp := []*big.Int{opts.GasPrice}
	fmt.Println(gp)
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice()
		if err != nil {
			return &common.Hash{}, err
		}
	}
	fmt.Println(gp)

	tx, err := t.TxFabric(n.Uint64(), to, opts.Value, opts.GasLimit, gp, data)
	if err != nil {
		return &common.Hash{}, err
	}
	_, err = t.client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return &common.Hash{}, err
	}
	log.Debug().Msgf("hash: %v from: %s", tx.Hash(), t.client.From())
	_, err = t.client.WaitAndReturnTxReceipt(tx.Hash())
	if err != nil {
		return &common.Hash{}, err
	}
	err = t.client.UnsafeIncreaseNonce()
	if err != nil {
		return &common.Hash{}, err
	}
	h := tx.Hash()
	return &h, nil
}
