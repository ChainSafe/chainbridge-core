package signAndSend

import (
	"context"
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/client"
	"github.com/ChainSafe/sygma-core/chains/evm/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/transactor/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type GasPricer interface {
	GasPrice(priority *uint8) ([]*big.Int, error)
}

type signAndSendTransactor struct {
	TxFabric       transaction.TxFabric
	gasPriceClient GasPricer
	client         client.Client
}

func NewSignAndSendTransactor(txFabric transaction.TxFabric, gasPriceClient GasPricer, client client.Client) transactor.Transactor {
	return &signAndSendTransactor{
		TxFabric:       txFabric,
		gasPriceClient: gasPriceClient,
		client:         client,
	}
}

func (t *signAndSendTransactor) Transact(to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error) {
	t.client.LockNonce()
	n, err := t.client.UnsafeNonce()
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, err
	}

	err = transactor.MergeTransactionOptions(&opts, &transactor.DefaultTransactionOptions)
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, err
	}

	gp := []*big.Int{opts.GasPrice}
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice(&opts.Priority)
		if err != nil {
			t.client.UnlockNonce()
			return &common.Hash{}, err
		}
	}

	tx, err := t.TxFabric(n.Uint64(), to, opts.Value, opts.GasLimit, gp, data)
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, err
	}

	h, err := t.client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		t.client.UnlockNonce()
		log.Error().Err(err)
		return &common.Hash{}, err
	}

	err = t.client.UnsafeIncreaseNonce()
	t.client.UnlockNonce()
	if err != nil {
		return &common.Hash{}, err
	}

	_, err = t.client.WaitAndReturnTxReceipt(h)
	if err != nil {
		return &common.Hash{}, err
	}

	return &h, nil
}
