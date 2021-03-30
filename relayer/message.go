package relayer

// Message is used as a generic format to communicate between chains
type Message struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	Type         string // type of bridge transfer
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   [32]byte
	Payload      []interface{} // data associated with event sequence
}

type Messager interface {
	GetSource() uint8
	GetDestination() uint8
	GetType() string
	GetDepositNonce() uint64
	GetResourceID() [32]byte
	GetPayload() []interface{} // Maybe this should be some bytes encoding
}
