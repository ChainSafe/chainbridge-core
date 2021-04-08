package relayer

import (
	"github.com/ethereum/go-ethereum/common"
)

// XCMessage is used as a generic format cross-chain communications
type XCMessager interface {
	GetSource() uint8
	GetDestination() uint8
	GetType() string
	GetDepositNonce() uint64
	GetResourceID() [32]byte
	GetPayload() []interface{} // Maybe this should be some bytes encoding
}

type ProposalStatus uint8

const (
	ProposalNotPassedStatus   ProposalStatus = 1
	ProposalStatusPassed      ProposalStatus = 2
	ProposalStatusTransferred ProposalStatus = 3
	ProposalStatusCancelled   ProposalStatus = 4
)

type Proposal interface {
	XCMessager
	GetProposalData() ([]byte, error)
	GetProposalDataHash(data []byte) common.Hash
	GetProposalStatus() ProposalStatus
	ShouldVoteFor() bool
}

type ProposalCreatorFn func(msg XCMessager) (Proposal, error)
