// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

import (
	"context"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/rs/zerolog/log"
)

type Metrics interface {
	TrackDepositMessage(m *message.Message)
}

type RelayedChain interface {
	PollEvents(ctx context.Context, sysErr chan<- error, msgChan chan []*message.Message)
	Write(messages []*message.Message)
	DomainID() uint8
}

func NewRelayer(chains []RelayedChain, metrics Metrics, messageProcessors ...message.MessageProcessor) *Relayer {
	return &Relayer{relayedChains: chains, messageProcessors: messageProcessors, metrics: metrics}
}

type Relayer struct {
	metrics           Metrics
	relayedChains     []RelayedChain
	registry          map[uint8]RelayedChain
	messageProcessors []message.MessageProcessor
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

// Route function winds destination writer by mapping DestinationID from message to registered writer.
func (r *Relayer) route(msgs []*message.Message) {
	destChain, ok := r.registry[msgs[0].Destination]
	if !ok {
		log.Error().Msgf("no resolver for destID %v to send message registered", msgs[0].Destination)
		return
	}

	for _, m := range msgs {
		r.metrics.TrackDepositMessage(m)

		for _, mp := range r.messageProcessors {
			if err := mp(m); err != nil {
				log.Error().Err(fmt.Errorf("error %w processing mesage %v", err, m))
				return
			}
		}
	}

	log.Debug().Msgf("Sending messages %+v to destination %v", msgs, destChain.DomainID())
	destChain.Write(msgs)
}

func (r *Relayer) addRelayedChain(c RelayedChain) {
	if r.registry == nil {
		r.registry = make(map[uint8]RelayedChain)
	}
	domainID := c.DomainID()
	r.registry[domainID] = c
}
