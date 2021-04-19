package relayer

import "github.com/ethereum/go-ethereum/common"

// XCMessage is used as a generic format cross-chain communications
type XCMessager interface {
	GetSource() uint8
	GetDestination() uint8
	GetDepositNonce() uint64
	GetResourceID() [32]byte
	GetPayload() []interface{} // Maybe this should be some bytes encoding
	String() string
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
	GetProposalData() []byte
	GetProposalDataHash() common.Hash
	ShouldBeVotedFor() bool
	ProposalIsComplete() bool
	ProposalIsReadyForExecute() bool
}
