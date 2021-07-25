package cliutils

const (
	Deposit       EventSig = "Deposit(uint8,bytes32,uint64)"
	ProposalEvent EventSig = "ProposalEvent(uint8,uint64,uint8,bytes32,bytes32)"
)

type ProposalStatus int

const (
	Inactive ProposalStatus = iota
	Active
	Passed
	Executed
	Cancelled
)
