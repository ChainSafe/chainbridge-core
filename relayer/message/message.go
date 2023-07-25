// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package message

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/types"
)

type TransferType string

type Metadata struct {
	Priority uint8
	Data     map[string]interface{}
}

const (
	FungibleTransfer    TransferType = "FungibleTransfer"
	NonFungibleTransfer TransferType = "NonFungibleTransfer"
	GenericTransfer     TransferType = "GenericTransfer"
)

type ProposalStatus struct {
	Status        uint8
	YesVotes      *big.Int
	YesVotesTotal uint8
	ProposedBlock *big.Int
}

const (
	ProposalStatusInactive uint8 = iota
	ProposalStatusActive
	ProposalStatusPassed // Ready to be executed
	ProposalStatusExecuted
	ProposalStatusCanceled
)

var (
	StatusMap = map[uint8]string{ProposalStatusInactive: "inactive", ProposalStatusActive: "active", ProposalStatusPassed: "passed", ProposalStatusExecuted: "executed", ProposalStatusCanceled: "canceled"}
)

type Message struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   types.ResourceID
	Payload      []interface{} // data associated with event sequence
	Metadata     Metadata      // Arbitrary data that will be most likely be used by the relayer
	Type         TransferType
}

func (m *Message) payloadToString() string {
	var payload bytes.Buffer
	for _, v := range m.Payload {
		bv, ok := v.([]byte)
		// All the payload should be in []byte type so this might never be a problem
		if !ok {
			continue
		}
		payload.Write(bv)
	}
	return payload.String()
}

func (m *Message) String() string {
	return fmt.Sprintf(
		`Source: %d; 
				Destination: %d, 
				DepositNonce: %d, 
				ResourceID: %x, 
				Payload: %s, 
				Metadata: %+v, 
				Type: %s`,
		m.Source, m.Destination, m.DepositNonce, m.ResourceId, m.PayloadToString(), m.Metadata, m.Type)
}

func NewMessage(
	source uint8,
	destination uint8,
	depositNonce uint64,
	resourceId types.ResourceID,
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
	return strconv.FormatInt(int64(m.Source), 10) + "-" + strconv.FormatInt(int64(m.DepositNonce), 10)
}
