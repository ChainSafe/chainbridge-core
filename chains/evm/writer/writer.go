package writer

import (
	"errors"
	"time"

	"github.com/ChainSafe/chainbridgev2/relayer"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Number of blocks to wait for an finalization event
const ExecuteBlockWatchLimit = 100

var BlockRetryInterval = time.Second * 5

type ProposalExecutor interface {
	ExecuteProposal(proposal relayer.Proposal)
}

type ProposalVoter interface {
	VoteProposal(proposal relayer.Proposal)
}

type VoterExecutor interface {
	ProposalExecutor
	ProposalVoter
}

type ProposalHandler func(msg relayer.XCMessager, handlerAddr string) (relayer.Proposal, error)
type ProposalHandlers map[ethcommon.Address]ProposalHandler

type BridgeReader interface {
	MatchResourceIDToHandlerAddress(rID [32]byte) (string, error)
}

type Writer struct {
	stop                  <-chan struct{}
	errChn                chan<- error
	handlers              ProposalHandlers
	bridgeReader          BridgeReader
	proposalVoterExecutor VoterExecutor
}

func NewWriter(stop <-chan struct{}, errChn chan<- error, ve VoterExecutor, bridgeReader BridgeReader) *Writer {
	return &Writer{
		stop:                  stop,
		errChn:                errChn,
		proposalVoterExecutor: ve,
		bridgeReader:          bridgeReader,
	}
}

func (w *Writer) Write(m relayer.XCMessager) {
	// Matching resource ID with handler.
	addr, err := w.bridgeReader.MatchResourceIDToHandlerAddress(m.GetResourceID())
	// Based on handler that registered on BridgeContract
	propHandler, err := w.MatchAddressWithHandlerFunc(addr)
	if err != nil {
		w.errChn <- err
		return
	}
	prop, err := propHandler(m, addr)
	if err != nil {
		w.errChn <- err
		return
	}

	if !prop.ShouldBeVotedFor() {
		if prop.ProposalIsReadyForExecute() {
			// We should not vote for this proposal but it is ready to be executed
			w.proposalVoterExecutor.ExecuteProposal(prop)
			return
		} else {
			return
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
				return
			}
			continue
		case <-w.stop:
			return

		}
	}
}

func (w *Writer) MatchAddressWithHandlerFunc(addr string) (ProposalHandler, error) {
	h, ok := w.handlers[ethcommon.HexToAddress(addr)]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return h, nil
}

func (w *Writer) RegisterHandler(address string, handler ProposalHandler) {
	w.handlers[ethcommon.HexToAddress(address)] = handler
}
