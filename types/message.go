package types

import (
	"strconv"
)

type Metadata struct {
	Data map[string]interface{}
}
type TransferType string
type Message struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   ResourceID
	Payload      []interface{} // data associated with event sequence
	Metadata     Metadata      // Arbitrary data that will be most likely be used by the relayer
	Type         TransferType
}

func NewMessage(
	source uint8,
	destination uint8,
	depositNonce uint64,
	resourceId ResourceID,
	transferType TransferType,
	payload []interface{},
	metadata Metadata,
) *Message {
	return &Message{
		source,
		destination,
		depositNonce,
		resourceId,
		payload,
		metadata,
		transferType,
	}
}

func (m Message) ID() string {
	return strconv.FormatInt(int64(m.Source), 10) + "-" + strconv.FormatInt(int64(m.Destination), 10) + "-" + strconv.FormatInt(int64(m.DepositNonce), 10)
}
