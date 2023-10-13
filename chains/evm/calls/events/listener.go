package events

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/observability"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span, logger := observability.CreateSpanAndLoggerFromContext(ctx, "relayer-core", "relayer.core.Listener.FetchDeposits", attribute.String("startBlock", startBlock.String()), attribute.String("endBlock", endBlock.String()))
	defer span.End()

	logs, err := l.client.FetchEventLogs(ctx, contractAddress, string(DepositSig), startBlock, endBlock)
	if err != nil {
		return nil, observability.LogAndRecordErrorWithStatus(nil, span, err, "failed FetchEventLogs")
	}
	deposits := make([]*Deposit, 0)

	for _, dl := range logs {
		d, err := l.UnpackDeposit(l.abi, dl.Data)
		if err != nil {
			_ = observability.LogAndRecordError(&logger, span, err, "failed unpacking deposit event log")
			continue
		}
		d.SenderAddress = common.BytesToAddress(dl.Topics[1].Bytes())
		deposits = append(deposits, d)
		observability.LogAndEvent(
			logger.Debug(),
			span,
			fmt.Sprintf("Found deposit log in block: %d, TxHash: %s, contractAddress: %s, sender: %s", dl.BlockNumber, dl.TxHash, dl.Address, d.SenderAddress),
			append(d.TraceEventAttributes(), attribute.String("tx.hash", dl.TxHash.Hex()))...)
	}
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
