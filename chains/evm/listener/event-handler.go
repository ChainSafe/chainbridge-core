package listener

import (
	"context"
	"fmt"
	"math/big"

	traceapi "go.opentelemetry.io/otel/trace"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/events"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
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
	ctxWithSpan, span := otel.Tracer("relayer-core").Start(ctx, "relayer.core.DepositEventHandler.HandleEvent")
	defer span.End()
	span.SetAttributes(attribute.String("startBlock", startBlock.String()), attribute.String("endBlock", endBlock.String()))
	logger := log.With().Str("trace_id", span.SpanContext().TraceID().String()).Logger()
	deposits, err := eh.eventListener.FetchDeposits(ctxWithSpan, eh.bridgeAddress, startBlock, endBlock)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("unable to fetch deposit events because of: %+v", err)
	}

	domainDeposits := make(map[uint8][]*message.Message)
	for _, d := range deposits {
		func(d *events.Deposit) {
			defer func() {
				if r := recover(); r != nil {
					log.Error().Err(err).Msgf("panic occured while handling deposit %+v", d)
				}
			}()

			m, err := eh.depositHandler.HandleDeposit(eh.domainID, d.DestinationDomainID, d.DepositNonce, d.ResourceID, d.Data, d.HandlerResponse)
			if err != nil {
				logger.Error().Err(err).Str("start block", startBlock.String()).Str("end block", endBlock.String()).Uint8("domainID", eh.domainID).Msgf("%v", err)
				span.SetStatus(codes.Error, err.Error())
				return
			}

			logger.Debug().Str("msg.id", m.ID()).Msgf("Resolved message %s in block range: %s-%s", m.String(), startBlock.String(), endBlock.String())
			span.AddEvent("Resolved message", traceapi.WithAttributes(attribute.String("msg.id", m.ID()), attribute.String("msg.type", string(m.Type))))
			domainDeposits[m.Destination] = append(domainDeposits[m.Destination], m)
		}(d)
	}

	for _, deposits := range domainDeposits {
		go func(d []*message.Message) {
			msgChan <- d
		}(deposits)
	}
	span.SetStatus(codes.Ok, "Deposits handled")
	return nil
}
