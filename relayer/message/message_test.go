// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package message

import (
	"math/big"
	"testing"
)

// TestExtractAmountTransferred tests extractAmountTransferred to extract the total amount
// transferred during the transfer event
func TestExtractAmountTransferred(t *testing.T) {
	// init instance of Message
	msg := &Message{
		Payload: []interface{}{
			big.NewInt(10).Bytes(), // 10 tokens
		},
	}

	// extract amount from message payload
	payloadAmount, err := msg.extractAmountTransferred()
	if err != nil {
		t.Fatalf("could not extract amount transferred: %v", err)
	}

	// init sample value to test against
	expectedAmount := 10.0

	if payloadAmount != expectedAmount {
		t.Fatal("amounts do not equal")
	}
}
