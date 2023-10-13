package signAndSend

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/observability"

	"go.opentelemetry.io/otel/attribute"
	traceapi "go.opentelemetry.io/otel/trace"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/common"
)

type signAndSendTransactor struct {
	TxFabric       calls.TxFabric
	gasPriceClient calls.GasPricer
	client         calls.ClientDispatcher
}

func NewSignAndSendTransactor(txFabric calls.TxFabric, gasPriceClient calls.GasPricer, client calls.ClientDispatcher) transactor.Transactor {
	return &signAndSendTransactor{
		TxFabric:       txFabric,
		gasPriceClient: gasPriceClient,
		client:         client,
	}
}

func (t *signAndSendTransactor) Transact(ctx context.Context, to *common.Address, data []byte, opts transactor.TransactOptions) (*common.Hash, error) {
	ctx, span, _ := observability.CreateSpanAndLoggerFromContext(ctx, "relayer-core", "relayer.core.Transactor.signAndSendTransactor.Transact")
	defer span.End()

	t.client.LockNonce()
	n, err := t.client.UnsafeNonce()
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "failed to call UnsafeNonce")
	}

	err = transactor.MergeTransactionOptions(&opts, &transactor.DefaultTransactionOptions)
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "failed to MergeTransactionOptions")
	}

	gp := []*big.Int{opts.GasPrice}
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice(&opts.Priority)
		if err != nil {
			t.client.UnlockNonce()
			return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "failed to define gas price")
		}
	}

	span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.StringSlice("tx.gp", calls.BigIntSliceToStringSlice(gp))))

	tx, err := t.TxFabric(n.Uint64(), to, opts.Value, opts.GasLimit, gp, data)
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "unable to call TxFabric")
	}

	h, err := t.client.SignAndSendTransaction(ctx, tx)
	if err != nil {
		t.client.UnlockNonce()
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "unable to SignAndSendTransaction")
	}

	span.AddEvent("Transaction sent", traceapi.WithAttributes(attribute.String("tx.hash", h.String())))

	err = t.client.UnsafeIncreaseNonce()
	t.client.UnlockNonce()
	if err != nil {
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "unable to UnsafeIncreaseNonce")
	}

	_, err = t.client.WaitAndReturnTxReceipt(h)
	if err != nil {
		return &common.Hash{}, observability.LogAndRecordErrorWithStatus(nil, span, err, "failed to WaitAndReturnTxReceipt")
	}
	return &h, nil
}
