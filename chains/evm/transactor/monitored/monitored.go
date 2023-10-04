package monitored

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/ChainSafe/sygma-core/chains/evm/client"
	"github.com/ChainSafe/sygma-core/chains/evm/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/transactor/transaction"
)

type GasPricer interface {
	GasPrice(priority *uint8) ([]*big.Int, error)
}

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
	txFabric       transaction.TxFabric
	gasPriceClient GasPricer
	client         client.Client

	maxGasPrice        *big.Int
	increasePercentage *big.Int

	pendingTxns map[common.Hash]RawTx
	txLock      sync.Mutex
}

// NewMonitoredTransactor creates an instance of a transactor
// that periodically checks sent transactions and resends them
// with higher gas if they are stuck.
//
// Gas price is increased by increasePercentage param which
// is a percentage value with which old gas price should be increased (e.g 15)
func NewMonitoredTransactor(
	txFabric transaction.TxFabric,
	gasPriceClient GasPricer,
	client client.Client,
	maxGasPrice *big.Int,
	increasePercentage *big.Int,
) *MonitoredTransactor {
	return &MonitoredTransactor{
		client:             client,
		gasPriceClient:     gasPriceClient,
		txFabric:           txFabric,
		pendingTxns:        make(map[common.Hash]RawTx),
		maxGasPrice:        maxGasPrice,
		increasePercentage: increasePercentage,
	}
}

func (t *MonitoredTransactor) Transact(to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error) {
	t.client.LockNonce()
	defer t.client.UnlockNonce()

	n, err := t.client.UnsafeNonce()
	if err != nil {
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
	tx, err := t.txFabric(rawTx.nonce, rawTx.to, rawTx.value, rawTx.gasLimit, rawTx.gasPrice, rawTx.data)
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

func (t *MonitoredTransactor) Monitor(
	ctx context.Context,
	resendInterval time.Duration,
	txTimeout time.Duration,
	tooNewTransaction time.Duration,
) {
	ticker := time.NewTicker(resendInterval)

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

				for oldHash, tx := range pendingTxCopy {
					receipt, err := t.client.TransactionReceipt(context.Background(), oldHash)
					if err == nil {
						if receipt.Status == types.ReceiptStatusSuccessful {
							log.Info().Uint64("nonce", tx.nonce).Msgf("Executed transaction %s with nonce %d", oldHash, tx.nonce)
						} else {
							log.Error().Uint64("nonce", tx.nonce).Msgf("Transaction %s failed on chain", oldHash)
						}

						delete(t.pendingTxns, oldHash)
						continue
					}

					if time.Since(tx.creationTime) > txTimeout {
						log.Error().Uint64("nonce", tx.nonce).Msgf("Transaction %s has timed out", oldHash)
						delete(t.pendingTxns, oldHash)
						continue
					}
					if time.Since(tx.submitTime) < tooNewTransaction {
						continue
					}

					hash, err := t.resendTransaction(&tx)
					if err != nil {
						log.Warn().Uint64("nonce", tx.nonce).Err(err).Msgf("Failed resending transaction %s", hash)
						continue
					}

					delete(t.pendingTxns, oldHash)
					t.pendingTxns[hash] = tx
				}
			}
		}
	}
}

func (t *MonitoredTransactor) resendTransaction(tx *RawTx) (common.Hash, error) {
	tx.gasPrice = t.IncreaseGas(tx.gasPrice)
	newTx, err := t.txFabric(tx.nonce, tx.to, tx.value, tx.gasLimit, tx.gasPrice, tx.data)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := t.client.SignAndSendTransaction(context.TODO(), newTx)
	if err != nil {
		return common.Hash{}, err
	}

	log.Debug().Uint64("nonce", tx.nonce).Msgf("Resent transaction with hash %s", hash)

	return hash, nil
}

// increase gas bumps gas price by preset percentage.
//
// If gas was 10 and the increaseFactor is 15 the new gas price
// would be 11 (it floors the value). In case the gas price didn't
// change it increases it by 1.
func (t *MonitoredTransactor) IncreaseGas(oldGp []*big.Int) []*big.Int {
	newGp := make([]*big.Int, len(oldGp))
	for i, gp := range oldGp {

		percentIncreaseValue := new(big.Int).Div(new(big.Int).Mul(gp, t.increasePercentage), big.NewInt(100))
		increasedGp := new(big.Int).Add(gp, percentIncreaseValue)
		if increasedGp.Cmp(t.maxGasPrice) != -1 {
			increasedGp = t.maxGasPrice
		}

		if gp.Cmp(increasedGp) == 0 {
			newGp[i] = new(big.Int).Add(gp, big.NewInt(1))
		} else {
			newGp[i] = increasedGp
		}
	}
	return newGp
}
