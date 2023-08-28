package events

import (
	"context"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	traceapi "go.opentelemetry.io/otel/trace"
)

type ChainClient interface {
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]ethTypes.Log, error)
}

type Listener struct {
	client ChainClient
	abi    abi.ABI
}

func NewListener(client ChainClient) *Listener {
	abi, _ := abi.JSON(strings.NewReader(consts.BridgeABI))
	return &Listener{
		client: client,
		abi:    abi,
	}
}

func (l *Listener) FetchDeposits(ctx context.Context, contractAddress common.Address, startBlock *big.Int, endBlock *big.Int) ([]*Deposit, error) {
	ctxWithSpan, span := otel.Tracer("relayer-core").Start(ctx, "relayer.core.Listener.FetchDeposits")
	defer span.End()
	span.SetAttributes(attribute.String("startBlock", startBlock.String()), attribute.String("endBlock", endBlock.String()))
	logger := log.With().Str("dd.trace_id", span.SpanContext().TraceID().String()).Logger()

	logs, err := l.client.FetchEventLogs(ctxWithSpan, contractAddress, string(DepositSig), startBlock, endBlock)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	deposits := make([]*Deposit, 0)

	for _, dl := range logs {
		d, err := l.UnpackDeposit(l.abi, dl.Data)
		if err != nil {
			logger.Error().Msgf("failed unpacking deposit event log: %v", err)
			span.RecordError(err)
			continue
		}

		d.SenderAddress = common.BytesToAddress(dl.Topics[1].Bytes())
		logger.Debug().Msgf("Found deposit log in block: %d, TxHash: %s, contractAddress: %s, sender: %s", dl.BlockNumber, dl.TxHash, dl.Address, d.SenderAddress)
		span.AddEvent("Found deposit", traceapi.WithAttributes(append(d.TraceEventAttributes(), attribute.String("tx.hash", dl.TxHash.Hex()))...))
		deposits = append(deposits, d)
	}
	span.SetStatus(codes.Ok, "Deposits fetched")
	return deposits, nil
}

func (l *Listener) UnpackDeposit(abi abi.ABI, data []byte) (*Deposit, error) {
	var dl Deposit

	err := abi.UnpackIntoInterface(&dl, "Deposit", data)
	if err != nil {
		return &Deposit{}, err
	}

	return &dl, nil
}
