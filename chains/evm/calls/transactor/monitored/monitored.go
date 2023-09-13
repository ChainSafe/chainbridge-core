package monitored

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridge-core/observability"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.opentelemetry.io/otel/attribute"
	traceapi "go.opentelemetry.io/otel/trace"
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

type RawTxWithTraceID struct {
	RawTx
	traceID traceapi.TraceID
}

type MonitoredTransactor struct {
	txFabric       calls.TxFabric
	gasPriceClient calls.GasPricer
	client         calls.ClientDispatcher

	maxGasPrice        *big.Int
	increasePercentage *big.Int

	pendingTxns map[common.Hash]RawTxWithTraceID

	txLock sync.Mutex
}

// NewMonitoredTransactor creates an instance of a transactor
// that periodically checks sent transactions and resends them
// with higher gas if they are stuck.
//
// Gas price is increased by increasePercentage param which
// is a percentage value with which old gas price should be increased (e.g 15)
func NewMonitoredTransactor(
	txFabric calls.TxFabric,
	gasPriceClient calls.GasPricer,
	client calls.ClientDispatcher,
	maxGasPrice *big.Int,
	increasePercentage *big.Int,
) *MonitoredTransactor {
	return &MonitoredTransactor{
		client:             client,
		gasPriceClient:     gasPriceClient,
		txFabric:           txFabric,
		pendingTxns:        make(map[common.Hash]RawTxWithTraceID),
		maxGasPrice:        maxGasPrice,
		increasePercentage: increasePercentage,
	}
}

func (t *MonitoredTransactor) Transact(ctx context.Context, to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error) {
	_, span, _ := observability.CreateSpanAndLoggerFromContext(ctx, "relayer-core", "relayer.core.evm.monitoredTransactor.Transact")
	defer span.End()

	t.client.LockNonce()
	defer t.client.UnlockNonce()

	n, err := t.client.UnsafeNonce()
	if err != nil {
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "failed to call UnsafeNonce")
	}

	err = transactor.MergeTransactionOptions(&opts, &transactor.DefaultTransactionOptions)
	if err != nil {
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "failed to MergeTransactionOptions")
	}

	gp := []*big.Int{opts.GasPrice}
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice(&opts.Priority)
		if err != nil {
			return &common.Hash{}, err
		}
	}
	span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.StringSlice("tx.gp", calls.BigIntSliceToStringSlice(gp))))

	rawTx := RawTxWithTraceID{
		RawTx{
			to:           to,
			nonce:        n.Uint64(),
			value:        opts.Value,
			gasLimit:     opts.GasLimit,
			gasPrice:     gp,
			data:         data,
			submitTime:   time.Now(),
			creationTime: time.Now(),
		},
		span.SpanContext().TraceID(),
	}
	tx, err := t.txFabric(rawTx.nonce, rawTx.to, rawTx.value, rawTx.gasLimit, rawTx.gasPrice, rawTx.data)
	if err != nil {
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "unable to call TxFabric")
	}

	h, err := t.client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "unable to SignAndSendTransaction")
	}
	span.AddEvent("Transaction sent", traceapi.WithAttributes(attribute.String("tx.hash", h.String())))

	t.txLock.Lock()
	t.pendingTxns[h] = rawTx
	t.txLock.Unlock()

	err = t.client.UnsafeIncreaseNonce()
	if err != nil {
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "unable to UnsafeIncreaseNonce")
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
				pendingTxCopy := make(map[common.Hash]RawTxWithTraceID, len(t.pendingTxns))
				for k, v := range t.pendingTxns {
					pendingTxCopy[k] = v
				}
				t.txLock.Unlock()

				for oldHash, tx := range pendingTxCopy {
					if time.Since(tx.submitTime) < tooNewTransaction {
						continue
					}
					txContextWithSpan, span, logger := observability.CreateSpanAndLoggerFromContext(
						traceapi.ContextWithSpanContext(ctx, traceapi.NewSpanContext(traceapi.SpanContextConfig{TraceID: tx.traceID})),
						"relayer-core",
						"relayer.core.evm.transactor.Monitor",
						attribute.String("tx.hash", oldHash.String()), attribute.Int64("tx.nonce", int64(tx.nonce)))
					receipt, err := t.client.TransactionReceipt(context.Background(), oldHash)
					if err == nil {
						if receipt.Status == types.ReceiptStatusSuccessful {
							observability.LogAndEvent(logger.Info(), span, fmt.Sprintf("Executed transaction %s with nonce %d", oldHash, tx.nonce))
						} else {
							_ = observability.LogAndRecordErrorWithStatus(&logger, span, fmt.Errorf("on-chain execution fail"), fmt.Sprintf("transaction %s failed on chain", oldHash))
						}
						span.End()
						delete(t.pendingTxns, oldHash)
						continue
					}

					if time.Since(tx.creationTime) > txTimeout {
						_ = observability.LogAndRecordErrorWithStatus(&logger, span, fmt.Errorf("transaction has timed out"), fmt.Sprintf("transaction %s failed on chain", oldHash))
						span.End()
						delete(t.pendingTxns, oldHash)
						continue
					}

					hash, err := t.resendTransaction(txContextWithSpan, &tx.RawTx)
					if err != nil {
						span.RecordError(fmt.Errorf("error resending transaction %w", err), traceapi.WithAttributes(attribute.String("tx.hash", oldHash.String()), attribute.Int64("tx.nonce", int64(tx.nonce))))
						logger.Warn().Uint64("nonce", tx.nonce).Err(err).Msgf("Failed resending transaction %s", oldHash)
						_ = observability.LogAndRecordError(&logger, span, err, "failed resending transaction")
						continue
					}
					span.AddEvent("Transaction resent", traceapi.WithAttributes(attribute.String("tx.newHash", hash.String())))
					span.End()

					delete(t.pendingTxns, oldHash)
					t.pendingTxns[hash] = tx
				}
			}
		}
	}
}

func (t *MonitoredTransactor) resendTransaction(ctx context.Context, tx *RawTx) (common.Hash, error) {
	ctx, span, logger := observability.CreateSpanAndLoggerFromContext(ctx, "relayer-core", "relayer.core.evm.transactor.Monitor.resendTransaction")
	defer span.End()
	tx.gasPrice = t.IncreaseGas(tx.gasPrice)
	if len(tx.gasPrice) > 1 {
		observability.LogAndEvent(logger.Debug(), span, "Calculated GasPrice", attribute.String("tx.gasTipCap", tx.gasPrice[0].String()), attribute.String("tx.gasFeeCap", tx.gasPrice[1].String()))
	} else {
		observability.LogAndEvent(logger.Debug(), span, "Calculated GasPrice", attribute.String("tx.gp", tx.gasPrice[0].String()))
	}
	newTx, err := t.txFabric(tx.nonce, tx.to, tx.value, tx.gasLimit, tx.gasPrice, tx.data)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := t.client.SignAndSendTransaction(ctx, newTx)
	if err != nil {
		return common.Hash{}, err
	}
	return hash, nil
}

// IncreaseGas bumps gas price by preset percentage.
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
