package listener

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/events"
	"github.com/ChainSafe/chainbridge-core/observability"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type EventListener interface {
	FetchDeposits(ctx context.Context, address common.Address, startBlock *big.Int, endBlock *big.Int) ([]*events.Deposit, error)
}

type DepositHandler interface {
	HandleDeposit(sourceID, destID uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*message.Message, error)
}

type DepositEventHandler struct {
	eventListener  EventListener
	depositHandler DepositHandler

	bridgeAddress common.Address
	domainID      uint8
}

func NewDepositEventHandler(eventListener EventListener, depositHandler DepositHandler, bridgeAddress common.Address, domainID uint8) *DepositEventHandler {
	return &DepositEventHandler{
		eventListener:  eventListener,
		depositHandler: depositHandler,
		bridgeAddress:  bridgeAddress,
		domainID:       domainID,
	}
}

func (eh *DepositEventHandler) HandleEvent(ctx context.Context, startBlock *big.Int, endBlock *big.Int, msgChan chan []*message.Message) error {
	ctxWithSpan, span, logger := observability.CreateSpanAndLoggerFromContext(
		ctx,
		"relayer-core",
		"relayer.core.DepositEventHandler.HandleEvent",
		attribute.String("startBlock", startBlock.String()), attribute.String("endBlock", endBlock.String()))
	defer span.End()

	deposits, err := eh.eventListener.FetchDeposits(ctxWithSpan, eh.bridgeAddress, startBlock, endBlock)
	if err != nil {
		return observability.LogAndRecordErrorWithStatus(nil, span, err, "unable to fetch deposit events")
	}

	domainDeposits := make(map[uint8][]*message.Message)
	for _, d := range deposits {
		func(d *events.Deposit) {
			defer func() {
				if r := recover(); r != nil {
					_ = observability.LogAndRecordError(&logger, span, errors.New("panic"), "panic occured while handling deposit", d.TraceEventAttributes()...)
				}
			}()

			m, err := eh.depositHandler.HandleDeposit(eh.domainID, d.DestinationDomainID, d.DepositNonce, d.ResourceID, d.Data, d.HandlerResponse)
			if err != nil {
				logger.Error().Err(err).Str("start block", startBlock.String()).Str("end block", endBlock.String()).Uint8("domainID", eh.domainID).Msgf("%v", err)
				span.SetStatus(codes.Error, err.Error())
				return
			}
			domainDeposits[m.Destination] = append(domainDeposits[m.Destination], m)
			observability.LogAndEvent(
				logger.Debug(),
				span,
				fmt.Sprintf("Resolved message %s in block range: %s-%s", m.String(), startBlock.String(), endBlock.String()),
				attribute.String("msg.id", m.ID()),
				attribute.String("msg.type", string(m.Type)))
		}(d)
	}

	for _, deposits := range domainDeposits {
		go func(d []*message.Message) {
			msgChan <- d
		}(deposits)
	}
	return nil
}
