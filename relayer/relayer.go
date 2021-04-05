package relayer

// Main chain abstraction
type Chainer interface {
	GetChainID() uint8
	GetListener() IListener
}

type Router interface {
	Send(dest uint8, m XCMessager) error
}

func NewRelayer(listeners []IListener, r Router) *Relayer {
	return &Relayer{listeners: listeners}
}

type Relayer struct {
	listeners []IListener
	router    Router
}

// Starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (r *Relayer) Start(stop <-chan struct{}, sysErr chan<- error) {
	for _, l := range r.listeners {
		go PollBlocks(l, stop, sysErr)
	}

	// PollBlocks should send all events to channel. Here we read events from channel and push them to the router.
	err := r.router.Send(dest, m)
}
