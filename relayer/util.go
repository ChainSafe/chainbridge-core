package relayer

import (
	"encoding/binary"
	"errors"
)

// extractAmount is a private function to extract and transform the transfer amount
// from the Payload field within the Message struct
func extractAmountTransferred(message *Message) (int, error) {
	// parse payload field from event log message to obtain transfer amount
	// payload slice of interfaces includes..
	// index 0: amount ([]byte)
	// index 1: destination recipient address ([]byte)

	b, ok := message.Payload[0].([]byte)
	if !ok {
		err := errors.New("could not cast interface to byte slice")
		return 0, err
	}

	// set payload amount by converting []byte => uint64
	payloadAmount := binary.BigEndian.Uint64(b)

	return int(payloadAmount), nil
}
