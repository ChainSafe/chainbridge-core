package monitored

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

type MonitoredTransactor struct {
	txFabric       calls.TxFabric
	gasPriceClient calls.GasPricer
	client         calls.ClientDispatcher

	maxGasPrice        *big.Int
	increasePercentage *big.Int

	pendingTxns      map[common.Hash]RawTx
	pendingTxnsTrace map[common.Hash]traceapi.TraceID

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
		pendingTxns:        make(map[common.Hash]RawTx),
		pendingTxnsTrace:   make(map[common.Hash]traceapi.TraceID),
		maxGasPrice:        maxGasPrice,
		increasePercentage: increasePercentage,
	}
}

func (t *MonitoredTransactor) Transact(ctx context.Context, to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error) {
	_, span := otel.Tracer("relayer-core").Start(ctx, "relayer.core.EVMListener.ListenToEvents")

	t.client.LockNonce()
	defer t.client.UnlockNonce()

	n, err := t.client.UnsafeNonce()
	if err != nil {
		span.RecordError(fmt.Errorf("unable to get unsafe nonce with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	err = transactor.MergeTransactionOptions(&opts, &transactor.DefaultTransactionOptions)
	if err != nil {
		span.RecordError(fmt.Errorf("unable to merge transaction options with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	gp := []*big.Int{opts.GasPrice}
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice(&opts.Priority)
		if err != nil {
			return &common.Hash{}, err
		}
	}
	if len(gp) > 1 {
		span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.String("tx.gasTipCap", gp[0].String()), attribute.String("tx.gasFeeCap", gp[1].String())))
	} else {
		span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.String("tx.gp", gp[0].String())))
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
		span.RecordError(fmt.Errorf("unable to call TxFabric with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	h, err := t.client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		span.RecordError(fmt.Errorf("unable to SignAndSendTransaction with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	t.txLock.Lock()
	t.pendingTxns[h] = rawTx
	t.pendingTxnsTrace[h] = span.SpanContext().TraceID()
	t.txLock.Unlock()

	err = t.client.UnsafeIncreaseNonce()
	if err != nil {
		span.RecordError(fmt.Errorf("unable to UnsafeIncreaseNonce with err: %w", err))
		span.End()
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
				pendingTxTraceIDCopy := make(map[common.Hash]traceapi.TraceID, len(t.pendingTxnsTrace))
				for k, v := range t.pendingTxnsTrace {
					pendingTxTraceIDCopy[k] = v
				}
				t.txLock.Unlock()

				for oldHash, tx := range pendingTxCopy {
					if time.Since(tx.submitTime) < tooNewTransaction {
						continue
					}
					tID, ok := pendingTxTraceIDCopy[oldHash]
					if ok {
						// Creating span context with existing TraceID
						spanCtx := traceapi.NewSpanContext(traceapi.SpanContextConfig{TraceID: tID, Remote: true})
						ctx = traceapi.ContextWithSpanContext(ctx, spanCtx)
					}
					ctx, span := otel.Tracer("relayer-sygma").Start(ctx, "relayer.sygma.evm.transactor.Monitor")
					logger := log.With().Str("dd.trace_id", span.SpanContext().TraceID().String()).Logger()

					receipt, err := t.client.TransactionReceipt(context.Background(), oldHash)
					if err == nil {
						if receipt.Status == types.ReceiptStatusSuccessful {
							logger.Info().Uint64("nonce", tx.nonce).Msgf("Executed transaction %s with nonce %d", oldHash, tx.nonce)
							span.AddEvent("Executed transaction", traceapi.WithAttributes(attribute.String("tx.hash", oldHash.String()), attribute.Int64("tx.nonce", int64(tx.nonce))))
							span.SetStatus(codes.Ok, "Executed transaction")
							span.End()
						} else {
							logger.Error().Uint64("nonce", tx.nonce).Msgf("Transaction %s failed on chain", oldHash)
							span.RecordError(fmt.Errorf("transaction execution failed on chain with error %w", err), traceapi.WithAttributes(attribute.String("tx.hash", oldHash.String()), attribute.Int64("tx.nonce", int64(tx.nonce))))
							span.SetStatus(codes.Error, "Transaction execution failed on chain")
							span.End()
						}
						delete(t.pendingTxns, oldHash)
						delete(t.pendingTxnsTrace, oldHash)
						continue
					}

					if time.Since(tx.creationTime) > txTimeout {
						logger.Error().Uint64("nonce", tx.nonce).Msgf("Transaction %s has timed out", oldHash)
						span.RecordError(fmt.Errorf("transaction has timed out"), traceapi.WithAttributes(attribute.String("tx.hash", oldHash.String()), attribute.Int64("tx.nonce", int64(tx.nonce))))
						span.End()
						delete(t.pendingTxns, oldHash)
						delete(t.pendingTxnsTrace, oldHash)
						continue
					}

					hash, err := t.resendTransaction(ctx, &tx)
					if err != nil {
						span.RecordError(fmt.Errorf("error resending transaction %w", err), traceapi.WithAttributes(attribute.String("tx.hash", oldHash.String()), attribute.Int64("tx.nonce", int64(tx.nonce))))
						logger.Warn().Uint64("nonce", tx.nonce).Err(err).Msgf("Failed resending transaction %s", oldHash)
						continue
					}
					span.AddEvent("Resending transaction", traceapi.WithAttributes(attribute.String("tx.newHash", hash.String())))
					span.End()

					delete(t.pendingTxns, oldHash)
					delete(t.pendingTxnsTrace, oldHash)
					t.pendingTxns[hash] = tx
					t.pendingTxnsTrace[hash] = tID
				}
			}
		}
	}
}

func (t *MonitoredTransactor) resendTransaction(ctx context.Context, tx *RawTx) (common.Hash, error) {
	tx.gasPrice = t.IncreaseGas(ctx, tx.gasPrice)
	newTx, err := t.txFabric(tx.nonce, tx.to, tx.value, tx.gasLimit, tx.gasPrice, tx.data)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := t.client.SignAndSendTransaction(ctx, newTx)
	if err != nil {
		return common.Hash{}, err
	}

	log.Debug().Uint64("nonce", tx.nonce).Msgf("Resent transaction with hash %s", hash)

	return hash, nil
}

// IncreaseGas bumps gas price by preset percentage.
//
// If gas was 10 and the increaseFactor is 15 the new gas price
// would be 11 (it floors the value). In case the gas price didn't
// change it increases it by 1.
func (t *MonitoredTransactor) IncreaseGas(ctx context.Context, oldGp []*big.Int) []*big.Int {
	_, span := otel.Tracer("relayer-core").Start(ctx, "relayer.sygma.evm.transactor.Monitor.IncreaseGas")
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
	if len(newGp) > 1 {
		span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.String("tx.gasTipCap", newGp[0].String()), attribute.String("tx.gasFeeCap", newGp[1].String())))
	} else {
		span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.String("tx.gp", newGp[0].String())))
	}
	span.End()
	return newGp
}
