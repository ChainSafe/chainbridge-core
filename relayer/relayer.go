package relayer

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func NewRelayer(listeners []IListener) *Relayer {
	return &Relayer{listeners: listeners}
}

type Relayer struct {
	listeners []IListener
	registry  map[uint8]ChainWriter
}

type ChainWriter interface {
	Write(XCMessager)
}

// Starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(stop <-chan struct{}, sysErr chan error) {
	messagesChannel := make(chan XCMessager)
	for _, l := range r.listeners {
		go PollEvents(l, stop, sysErr, messagesChannel)
	}
	select {
	case m := <-messagesChannel:
		go r.Send(m)
	case _ = <-stop:
		return
	}
}

func (r *Relayer) Send(m XCMessager) {
	w, ok := r.registry[m.GetDestination()]
	if !ok {
		log.Error().Msgf(fmt.Sprintf("no resolver for destID %v to send message registered", m.GetDestination()))
		return
	}
	log.Debug().Msgf("Sending message %+v to destination %v", m, m.GetDestination())
	w.Write(m)
}

func (r *Relayer) SetWriter(chainID uint8, w ChainWriter) {
	if r.registry == nil {
		r.registry = make(map[uint8]ChainWriter)
	}
	r.registry[chainID] = w
}
