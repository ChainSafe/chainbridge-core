package relayer

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
)

type BlockStorer interface {
	StoreBlock(block *big.Int, chainID uint8) error
	GetLastStoredBlock(chainID uint8) error
	Close() error
}

type RelayedChain interface {
	PollEvents(bs BlockStorer, stop <-chan struct{}, sysErr chan<- error, eventsChan chan XCMessager)
	Write(XCMessager)
	ChainID() uint8
}

func NewRelayer(chains []RelayedChain, bs BlockStorer) *Relayer {
	return &Relayer{relayedChains: chains, bs: bs}
}

type Relayer struct {
	relayedChains []RelayedChain
	registry      map[uint8]RelayedChain
	bs            BlockStorer
}

// Starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(stop <-chan struct{}, sysErr chan error) {
	messagesChannel := make(chan XCMessager)
	for _, c := range r.relayedChains {
		r.addRelayedChain(c)
		go c.PollEvents(r.bs, stop, sysErr, messagesChannel)
	}
	for {
		select {
		case m := <-messagesChannel:
			go r.Route(m)
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
