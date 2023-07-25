package listener

import (
	"context"
	"fmt"
	"math/big"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/events"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
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

func (eh *DepositEventHandler) HandleEvent(startBlock *big.Int, endBlock *big.Int, msgChan chan []*message.Message) error {
	tp := otel.GetTracerProvider()
	ctxWithSpan, span := tp.Tracer("relayer-core-tracer").Start(context.Background(), "relayer.core.DepositEventHandler.HandleEvent")
	defer span.End()
	span.SetAttributes(attribute.String("startBlock", startBlock.String()), attribute.String("endBlock", endBlock.String()))

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
				log.Error().Err(err).Str("start block", startBlock.String()).Str("end block", endBlock.String()).Uint8("domainID", eh.domainID).Msgf("%v", err)
				span.SetStatus(codes.Error, err.Error())
				return
			}

			log.Debug().Msgf("Resolved message %+v in block range: %s-%s", m, startBlock.String(), endBlock.String())
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
