package relayer

type TransferType string

var FungibleTransfer TransferType = "FungibleTransfer"
var NonFungibleTransfer TransferType = "NonFungibleTransfer"
var GenericTransfer TransferType = "GenericTransfer"

// XCMessage is used as a generic format cross-chain communications
type XCMessage struct {
	Source       uint8        // Source where message was initiated
	Destination  uint8        // Destination chain of message
	Type         TransferType // type of bridge transfer
	DepositNonce uint64       // Nonce for the deposit
	ResourceId   [32]byte
	Payload      []interface{} // data associated with event sequence
}

type XCMessager interface {
	GetSource() uint8
	GetDestination() uint8
	GetType() string
	GetDepositNonce() uint64
	GetResourceID() [32]byte
	GetPayload() []interface{} // Maybe this should be some bytes encoding
}
