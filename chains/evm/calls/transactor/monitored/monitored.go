package monitored

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
)

type RawTx struct {
	nonce        uint64
	to           *common.Address
	value        *big.Int
	gasLimit     uint64
	gasPrice     []*big.Int
	data         []byte
	submitTime   time.Time
	creationTime time.Time
}

type MonitoredTransactor struct {
	TxFabric       calls.TxFabric
	gasPriceClient calls.GasPricer
	client         calls.ClientDispatcher

	pendingTxns map[common.Hash]RawTx
	txLock      sync.Mutex

	resendInterval time.Duration
	txTimeout      time.Duration
	increaseFactor *big.Int
	maxGasPrice    *big.Int
}

func NewMonitoredTransactor() *MonitoredTransactor {
	t := &MonitoredTransactor{}
	return t
}

func (t *MonitoredTransactor) Transact(to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error) {
	t.client.LockNonce()
	defer t.client.UnlockNonce()

	n, err := t.client.UnsafeNonce()
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, err
	}

	err = transactor.MergeTransactionOptions(&opts, &transactor.DefaultTransactionOptions)
	if err != nil {
		return &common.Hash{}, err
	}

	gp := []*big.Int{opts.GasPrice}
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice(&opts.Priority)
		if err != nil {
			return &common.Hash{}, err
		}
	}

	rawTx := RawTx{
		to:           to,
		nonce:        n.Uint64(),
		value:        opts.Value,
		gasLimit:     opts.GasLimit,
		gasPrice:     gp,
		data:         data,
		submitTime:   time.Now(),
		creationTime: time.Now(),
	}
	tx, err := t.TxFabric(rawTx.nonce, rawTx.to, rawTx.value, rawTx.gasLimit, rawTx.gasPrice, rawTx.data)
	if err != nil {
		return &common.Hash{}, err
	}

	h, err := t.client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return &common.Hash{}, err
	}

	t.txLock.Lock()
	t.pendingTxns[h] = rawTx
	t.txLock.Unlock()

	err = t.client.UnsafeIncreaseNonce()
	if err != nil {
		return &common.Hash{}, err
	}

	return &h, nil
}

func (t *MonitoredTransactor) Monitor(ctx context.Context) {
	ticker := time.NewTicker(t.resendInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			{
				t.txLock.Lock()
				pendingTxCopy := make(map[common.Hash]RawTx, len(t.pendingTxns))
				for k, v := range t.pendingTxns {
					pendingTxCopy[k] = v
				}
				t.txLock.Unlock()

				for hash, tx := range pendingTxCopy {
					receipt, err := t.client.TransactionReceipt(context.Background(), hash)
					if err == nil {
						if receipt.Status == types.ReceiptStatusSuccessful {
							log.Info().Msgf("Executed transaction %s with nonce %d", hash, tx.nonce)
						} else {
							log.Error().Msgf("Transaction %s failed on chain with nonce %d", hash, tx.nonce)
						}

						delete(t.pendingTxns, hash)
						continue
					}

					if time.Since(tx.creationTime) > t.txTimeout {
						log.Error().Msgf("Transaction %s with nonce %d has timed out", hash, tx.nonce)
						delete(t.pendingTxns, hash)
						continue
					}
					// avoid resending transaction if it just submitted
					if time.Since(tx.submitTime) < time.Minute {
						continue
					}

					hash, err := t.resendTransaction(&tx)
					if err != nil {
						log.Error().Err(err).Msgf("Failed resending transaction %s with nonce %d", hash, tx.nonce)
						continue
					}

					delete(t.pendingTxns, hash)
					t.pendingTxns[hash] = tx
				}
			}
		}
	}
}

func (t *MonitoredTransactor) resendTransaction(tx *RawTx) (common.Hash, error) {
	tx.gasPrice = t.increaseGas(tx.gasPrice)
	newTx, err := t.TxFabric(tx.nonce, tx.to, tx.value, tx.gasLimit, tx.gasPrice, tx.data)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := t.client.SignAndSendTransaction(context.TODO(), newTx)
	if err != nil {
		return common.Hash{}, err
	}

	return hash, nil
}

func (t *MonitoredTransactor) increaseGas(oldGp []*big.Int) []*big.Int {
	newGp := make([]*big.Int, len(oldGp))
	for i, gp := range oldGp {
		percentIncreaseValue := new(big.Int).Div(new(big.Int).Mul(gp, t.increaseFactor), big.NewInt(100))
		increasedGp := new(big.Int).Add(gp, percentIncreaseValue)
		if increasedGp.Cmp(t.maxGasPrice) != -1 {
			increasedGp = t.maxGasPrice
		}
		newGp[i] = increasedGp
	}
	return newGp
}
