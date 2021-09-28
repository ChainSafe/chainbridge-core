// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

import (
	"fmt"
	"net/http"

	"github.com/ChainSafe/chainbridge-core/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type MessageProcessor func(message *Message) error

type RelayedChain interface {
	PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan *Message)
	Write(message *Message) error
	DomainID() uint8
}

func NewRelayer(chains []RelayedChain, messageProcessors ...MessageProcessor) *Relayer {
	return &Relayer{relayedChains: chains, messageProcessors: messageProcessors}
}

type Relayer struct {
	relayedChains     []RelayedChain
	registry          map[uint8]RelayedChain
	messageProcessors []MessageProcessor
}

// Starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(stop <-chan struct{}, sysErr chan error) {
	log.Debug().Msgf("Starting relayer")
	messagesChannel := make(chan *Message)

	// init new instance of ChainMetrics
	chainMetrics := metrics.NewChainMetrics()

	// init new mux router
	router := mux.NewRouter()

	// register path + handler
	router.Path("/metrics").Handler(promhttp.Handler())

	// start http server in non-blocking goroutine
	go func() {
		log.Fatal().Err(http.ListenAndServe(":2112", router))
	}()
	log.Debug().Msg("listening on: http://localhost:2112/metrics")

	for _, c := range r.relayedChains {
		log.Debug().Msgf("Starting chain %v", c.DomainID())
		r.addRelayedChain(c)
		go c.PollEvents(stop, sysErr, messagesChannel)
	}
	for {
		select {
		case m := <-messagesChannel:
			go r.route(m, chainMetrics)
			continue
		case _ = <-stop:
			return
		}
	}
}

// Route function winds destination writer by mapping DestinationID from message to registered writer.
func (r *Relayer) route(m *Message, chainMetrics *metrics.ChainMetrics) {
	destChain, ok := r.registry[m.Destination]
	if !ok {
		log.Error().Msgf("no resolver for destID %v to send message registered", m.Destination)
		return
	}

	// extract amount from Payload field
	// TODO: if the message is not for ERC20 transfer that panics or errors could appear
	payloadAmount, err := m.extractAmountTransferred()
	if err != nil {
		log.Error().Err(err)
		return
	}

	// increment chain metrics
	chainMetrics.AmountTransferred.Add(payloadAmount)
	chainMetrics.NumberOfTransfers.Inc()

	for _, mp := range r.messageProcessors {
		if err := mp(m); err != nil {
			log.Error().Err(fmt.Errorf("error %w processing mesage %v", err, m))
			return
		}
	}

	log.Debug().Msgf("Sending message %+v to destination %v", m, m.Destination)
	if err := destChain.Write(m); err != nil {
		log.Error().Err(err).Msgf("writing message %+v", m)
		return
	}
}

func (r *Relayer) addRelayedChain(c RelayedChain) {
	if r.registry == nil {
		r.registry = make(map[uint8]RelayedChain)
	}
	domainID := c.DomainID()
	r.registry[domainID] = c
}
