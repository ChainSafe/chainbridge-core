package relayer

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type RelayedChain interface {
	PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan XCMessager)
	Write(XCMessager)
	ChainID() uint8
}

func NewRelayer(chains []RelayedChain) *Relayer {
	return &Relayer{relayedChains: chains}
}

type Relayer struct {
	relayedChains []RelayedChain
	registry      map[uint8]RelayedChain
}

// Starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(stop <-chan struct{}, sysErr chan error) {
	log.Debug().Msgf("Starting relayer")
	messagesChannel := make(chan XCMessager)
	for _, c := range r.relayedChains {
		log.Debug().Msgf("Starting chain %v", c.ChainID())

		r.addRelayedChain(c)
		go c.PollEvents(stop, sysErr, messagesChannel)
	}
	for {
		select {
		case m := <-messagesChannel:
			go r.Route(m)
			continue
		case _ = <-stop:
			return
		}
	}
}

// Route function winds destination writer by mapping DestinationID from message to registered writer.
func (r *Relayer) Route(m XCMessager) {
	w, ok := r.registry[m.GetDestination()]
	if !ok {
		log.Error().Msgf(fmt.Sprintf("no resolver for destID %v to send message registered", m.GetDestination()))
		return
	}
	log.Debug().Msgf("Sending message %+v to destination %v", m, m.GetDestination())
	w.Write(m)
}

func (r *Relayer) addRelayedChain(c RelayedChain) {
	if r.registry == nil {
		r.registry = make(map[uint8]RelayedChain)
	}
	chainID := c.ChainID()
	r.registry[chainID] = c
}
