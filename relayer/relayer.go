package relayer

// Main chain abstraction
type Chainer interface {
	GetChainID() uint8
	GetListener() IListener
}

func NewRelayer(listeners []IListener) *Relayer {
	return &Relayer{listeners: listeners}
}

type Relayer struct {
	listeners []IListener
}

// Starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(stop <-chan struct{}, sysErr chan<- error) {
	for _, l := range r.listeners {
		go PollBlocks(l, stop, sysErr)
	}
}

//type Registry struct {
//	registry map[uint8]IWriter
//}
//
//// Add middleware for registry.
//func (r *Registry) SetWriter(id uint8, w IWriter) {
//	if r.registry == nil {
//		r.registry = make(map[uint8]IWriter)
//	}
//	r.registry[id] = w
//}
//
//func (r *Registry) Send(dest uint8, m XCMessager) error {
//	log.Debug().Msgf("Sending message %+v to destination %v", m, dest)
//	return nil
//}
