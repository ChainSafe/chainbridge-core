package writer

import (
	"errors"
	"time"

	"github.com/ChainSafe/chainbridgev2/relayer"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

var BlockRetryInterval = time.Second * 5

type VoterExecutor interface {
	ExecuteProposal(proposal relayer.Proposal)
	VoteProposal(proposal relayer.Proposal)
	MatchResourceIDToHandlerAddress(rID [32]byte) (string, error)
}

type ProposalHandler func(msg relayer.XCMessager, handlerAddr string) (relayer.Proposal, error)
type ProposalHandlers map[ethcommon.Address]ProposalHandler

type EVMWriter struct {
	stop                  <-chan struct{}
	handlers              ProposalHandlers
	proposalVoterExecutor VoterExecutor
}

func NewWriter(ve VoterExecutor) *EVMWriter {
	return &EVMWriter{
		proposalVoterExecutor: ve,
		handlers:              make(map[ethcommon.Address]ProposalHandler),
	}
}

func (w *EVMWriter) Write(m relayer.XCMessager) error {
	// Matching resource ID with handler.
	addr, err := w.proposalVoterExecutor.MatchResourceIDToHandlerAddress(m.GetResourceID())
	// Based on handler that registered on BridgeContract
	propHandler, err := w.MatchAddressWithHandlerFunc(addr)
	if err != nil {
		return err
	}
	prop, err := propHandler(m, addr)
	if err != nil {
		return err
	}

	if !prop.ShouldBeVotedFor() {
		if prop.ProposalIsReadyForExecute() {
			// We should not vote for this proposal but it is ready to be executed
			w.proposalVoterExecutor.ExecuteProposal(prop)
			return nil
		} else {
			return nil
		}
	}
	w.proposalVoterExecutor.VoteProposal(prop)
	// Checking every 5 seconds does proposal is ready to be executed
	// TODO: update infinity loop to break after some period of time
	for {
		select {
		case <-time.After(BlockRetryInterval):
			if prop.ProposalIsReadyForExecute() {
				w.proposalVoterExecutor.ExecuteProposal(prop)
				return nil
			}
			continue
		case <-w.stop:
			return nil

		}
	}
}

func (w *EVMWriter) MatchAddressWithHandlerFunc(addr string) (ProposalHandler, error) {
	h, ok := w.handlers[ethcommon.HexToAddress(addr)]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return h, nil
}

func (w *EVMWriter) RegisterProposalHandler(address string, handler ProposalHandler) {
	w.handlers[ethcommon.HexToAddress(address)] = handler
}
