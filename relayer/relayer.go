// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

import (
	"context"

	"github.com/ChainSafe/sygma-core/types"
	"github.com/rs/zerolog/log"
)

type DepositMeter interface {
	TrackDepositMessage(m *types.Message)
	TrackExecutionError(m *types.Message)
	TrackSuccessfulExecutionLatency(m *types.Message)
}

type RelayedChain interface {
	PollEvents(ctx context.Context, sysErr chan<- error)
	Write(messages []*types.Message) error
	DomainID() uint8
}

func NewRelayer(chains []RelayedChain, metrics DepositMeter) *Relayer {
	return &Relayer{relayedChains: chains, metrics: metrics}
}

type Relayer struct {
	metrics       DepositMeter
	relayedChains []RelayedChain
	registry      map[uint8]RelayedChain
}

// Start function starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(ctx context.Context, msgChan chan []*types.Message, sysErr chan error) {
	log.Debug().Msgf("Starting relayer")

	for _, c := range r.relayedChains {
		log.Debug().Msgf("Starting chain %v", c.DomainID())
		r.addRelayedChain(c)
		go c.PollEvents(ctx, sysErr)
	}

	for {
		select {
		case m := <-msgChan:
			go r.route(m)
			continue
		case <-ctx.Done():
			return
		}
	}
}

// Route function runs destination writer by mapping DestinationID from message to registered writer.
func (r *Relayer) route(msgs []*types.Message) {
	destChain, ok := r.registry[msgs[0].Destination]
	if !ok {
		log.Error().Msgf("no resolver for destID %v to send message registered", msgs[0].Destination)
		return
	}

	log.Debug().Msgf("Sending messages %+v to destination %v", msgs, destChain.DomainID())
	err := destChain.Write(msgs)
	if err != nil {
		for _, m := range msgs {
			log.Err(err).Msgf("Failed sending messages %+v to destination %v", m, destChain.DomainID())
			r.metrics.TrackExecutionError(m)
		}
		return
	}

	for _, m := range msgs {
		r.metrics.TrackSuccessfulExecutionLatency(m)
	}
}

func (r *Relayer) addRelayedChain(c RelayedChain) {
	if r.registry == nil {
		r.registry = make(map[uint8]RelayedChain)
	}
	domainID := c.DomainID()
	r.registry[domainID] = c
}
