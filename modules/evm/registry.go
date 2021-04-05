package evm

import (
	"errors"
	"fmt"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/rs/zerolog/log"
)

type Registry struct {
	registry map[uint8]relayer.ChainWriter
}

// Add middleware for registry.
func (r *Registry) SetWriter(chainID uint8, w relayer.ChainWriter) {
	if r.registry == nil {
		r.registry = make(map[uint8]relayer.ChainWriter)
	}
	r.registry[chainID] = w
}

func (r *Registry) Send(dest uint8, m relayer.XCMessager) error {
	w, ok := r.registry[dest]
	if !ok {
		return errors.New(fmt.Sprintf("no resolver for destID %v to send message registered", dest))
	}
	log.Debug().Msgf("Sending message %+v to destination %v", m, dest)
	w.Write(m)
	return nil
}
