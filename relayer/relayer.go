// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

import (
	"context"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/relayer/message/processors"

	"go.opentelemetry.io/otel/codes"

	"github.com/ChainSafe/chainbridge-core/observability"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
)

type DepositMeter interface {
	TrackDepositMessage(m *message.Message)
	TrackExecutionError(m *message.Message)
	TrackSuccessfulExecutionLatency(m *message.Message)
}

type RelayedChain interface {
	PollEvents(ctx context.Context, sysErr chan<- error, msgChan chan []*message.Message)
	Write(ctx context.Context, messages []*message.Message) error
	DomainID() uint8
}

func NewRelayer(chains []RelayedChain, metrics DepositMeter, messageProcessors ...processors.MessageProcessor) *Relayer {
	return &Relayer{relayedChains: chains, messageProcessors: messageProcessors, metrics: metrics}
}

type Relayer struct {
	metrics           DepositMeter
	relayedChains     []RelayedChain
	registry          map[uint8]RelayedChain
	messageProcessors []processors.MessageProcessor
}

// Start function starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(ctx context.Context, sysErr chan error) {
	log.Debug().Msgf("Starting relayer")
	messagesChannel := make(chan []*message.Message)
	for _, c := range r.relayedChains {
		log.Debug().Msgf("Starting chain %v", c.DomainID())
		r.addRelayedChain(c)
		go c.PollEvents(ctx, sysErr, messagesChannel)
	}
	for {
		select {
		case m := <-messagesChannel:
			go r.route(m)
			continue
		case <-ctx.Done():
			return
		}
	}
}

// Route function runs destination writer by mapping DestinationID from message to registered writer.
func (r *Relayer) route(msgs []*message.Message) {
	ctx, span, logger := observability.CreateSpanAndLoggerFromContext(context.Background(), "relayer-core", "relayer.core.Route")
	defer span.End()

	destChain, ok := r.registry[msgs[0].Destination]
	if !ok {
		_ = observability.LogAndRecordErrorWithStatus(&logger, span, fmt.Errorf("no resolver for destID %v to send message registered", msgs[0].Destination), "Routing failed")
		return
	}
	for _, m := range msgs {
		observability.LogAndEvent(
			logger.Info(),
			span,
			fmt.Sprintf("routing message %s", m.String()),
			attribute.String("msg.id", m.ID()),
			attribute.String("msg.type", string(m.Type)))
		r.metrics.TrackDepositMessage(m)
		for _, mp := range r.messageProcessors {
			if err := mp(ctx, m); err != nil {
				_ = observability.LogAndRecordErrorWithStatus(&logger, span, err, "message processing fail", attribute.String("msg.id", m.ID()))
				return
			}
		}
	}

	err := destChain.Write(ctx, msgs)
	if err != nil {
		for _, m := range msgs {
			_ = observability.LogAndRecordErrorWithStatus(&logger, span, err, "failed sending message to destination", attribute.String("msg.id", m.ID()))
			r.metrics.TrackExecutionError(m)
		}
		return
	}
	for _, m := range msgs {
		r.metrics.TrackSuccessfulExecutionLatency(m)
	}
	span.SetStatus(codes.Ok, "messages routed")
}

func (r *Relayer) addRelayedChain(c RelayedChain) {
	if r.registry == nil {
		r.registry = make(map[uint8]RelayedChain)
	}
	domainID := c.DomainID()
	r.registry[domainID] = c
}
