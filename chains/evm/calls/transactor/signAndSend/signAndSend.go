package signAndSend

import (
	"context"
	"fmt"
	"math/big"

	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("relayer-core").Start(ctx, "relayer.core.Transactor.signAndSendTransactor.Transact")
	t.client.LockNonce()
	n, err := t.client.UnsafeNonce()
	if err != nil {
		t.client.UnlockNonce()
		span.RecordError(fmt.Errorf("unable to get unsafe nonce with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	err = transactor.MergeTransactionOptions(&opts, &transactor.DefaultTransactionOptions)
	if err != nil {
		t.client.UnlockNonce()
		span.RecordError(fmt.Errorf("unable to merge transaction options with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	gp := []*big.Int{opts.GasPrice}
	if opts.GasPrice.Cmp(big.NewInt(0)) == 0 {
		gp, err = t.gasPriceClient.GasPrice(&opts.Priority)
		if err != nil {
			t.client.UnlockNonce()
			span.RecordError(fmt.Errorf("unable to define gas price with err: %w", err))
			span.End()
			return &common.Hash{}, err
		}
	}

	if len(gp) > 1 {
		span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.String("tx.gasTipCap", gp[0].String()), attribute.String("tx.gasFeeCap", gp[1].String())))
	} else {
		span.AddEvent("Calculated GasPrice", traceapi.WithAttributes(attribute.String("tx.gp", gp[0].String())))
	}

	tx, err := t.TxFabric(n.Uint64(), to, opts.Value, opts.GasLimit, gp, data)
	if err != nil {
		t.client.UnlockNonce()
		span.RecordError(fmt.Errorf("unable to call TxFabric with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	h, err := t.client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		t.client.UnlockNonce()
		span.RecordError(fmt.Errorf("unable to SignAndSendTransaction with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	err = t.client.UnsafeIncreaseNonce()
	t.client.UnlockNonce()
	if err != nil {
		span.RecordError(fmt.Errorf("unable to UnsafeIncreaseNonce with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}

	_, err = t.client.WaitAndReturnTxReceipt(h)
	if err != nil {
		span.RecordError(fmt.Errorf("unable to WaitAndReturnTxReceipt with err: %w", err))
		span.End()
		return &common.Hash{}, err
	}
	span.SetStatus(codes.Ok, "Transaction sent")
	span.End()
	return &h, nil
}
