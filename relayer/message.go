// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package relayer

type TransferType string

const (
	FungibleTransfer    TransferType = "FungibleTransfer"
	NonFungibleTransfer TransferType = "NonFungibleTransfer"
	GenericTransfer     TransferType = "GenericTransfer"
)

type ProposalStatus uint8

const (
	ProposalStatusInactive ProposalStatus = 0
	ProposalStatusActive   ProposalStatus = 1
	ProposalStatusPassed   ProposalStatus = 2 // Ready to be executed
	ProposalStatusExecuted ProposalStatus = 3
	ProposalStatusCanceled ProposalStatus = 4
)

var (
	StatusMap = map[ProposalStatus]string{ProposalStatusInactive: "inactive", ProposalStatusActive: "active", ProposalStatusPassed: "passed", ProposalStatusExecuted: "executed", ProposalStatusCanceled: "canceled"}
)

type Message struct {
	Source       uint8  // Source where message was initiated
	Destination  uint8  // Destination chain of message
	DepositNonce uint64 // Nonce for the deposit
	ResourceId   [32]byte
	Payload      []interface{} // data associated with event sequence
	Type         TransferType
}
