package writer

import (
	"errors"
	"time"

	"github.com/ChainSafe/chainbridgev2/chains/evm"

	"github.com/ChainSafe/chainbridgev2/relayer"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var BlockRetryInterval = time.Second * 5

type VoterExecutor interface {
	ExecuteProposal(bridgeAddress string, proposal *evm.Proposal) error
	VoteProposal(bridgeAddress string, proposal *evm.Proposal) error
	MatchResourceIDToHandlerAddress(bridgeAddress string, rID [32]byte) (string, error)
	ProposalStatus(bridgeAddress string, proposal *evm.Proposal) (relayer.ProposalStatus, error)
	VotedBy(bridgeAddress string, p *evm.Proposal) bool
}

type ProposalHandler func(msg *relayer.Message, handlerAddr string) (*evm.Proposal, error)
type ProposalHandlers map[ethcommon.Address]ProposalHandler

type EVMVoter struct {
	stop                  <-chan struct{}
	handlers              ProposalHandlers
	proposalVoterExecutor VoterExecutor
}

func NewWriter(ve VoterExecutor) *EVMVoter {
	return &EVMVoter{
		proposalVoterExecutor: ve,
		handlers:              make(map[ethcommon.Address]ProposalHandler),
	}
}

func (w *EVMVoter) VoteProposal(m *relayer.Message, bridgeAddress string) error {
	// Matching resource ID with handler.
	addr, err := w.proposalVoterExecutor.MatchResourceIDToHandlerAddress(bridgeAddress, m.ResourceId)
	// Based on handler that registered on BridgeContract
	handleProposal, err := w.MatchAddressWithHandlerFunc(addr)
	if err != nil {
		return err
	}
	prop, err := handleProposal(m, addr)
	if err != nil {
		return err
	}

	ps, err := w.proposalVoterExecutor.ProposalStatus(bridgeAddress, prop)
	if err != nil {
		log.Error().Err(err).Msgf("error getting proposal status %+v", prop)
	}

	votedByCurrentExecutor := w.proposalVoterExecutor.VotedBy(bridgeAddress, prop)

	if votedByCurrentExecutor || ps == relayer.ProposalStatusPassed || ps == relayer.ProposalStatusCanceled || ps == relayer.ProposalStatusExecuted {
		if ps == relayer.ProposalStatusPassed {
			// We should not vote for this proposal but it is ready to be executed
			err = w.proposalVoterExecutor.ExecuteProposal(bridgeAddress, prop)
			if err != nil {
				return err
			}
			return nil
		} else {
			return nil
		}
	}
	err = w.proposalVoterExecutor.VoteProposal(bridgeAddress, prop)
	if err != nil {
		return err
	}
	// Checking every 5 seconds does proposal is ready to be executed
	// TODO: somehow update infinity loop to break after some period of time
	for {
		select {
		case <-time.After(BlockRetryInterval):
			ps, err := w.proposalVoterExecutor.ProposalStatus(bridgeAddress, prop)
			if err != nil {
				log.Error().Err(err).Msgf("error getting proposal status %+v", prop)
				return err
			}
			if ps == relayer.ProposalStatusPassed {
				err = w.proposalVoterExecutor.ExecuteProposal(bridgeAddress, prop)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case <-w.stop:
			return nil

		}
	}
}

func (w *EVMVoter) MatchAddressWithHandlerFunc(addr string) (ProposalHandler, error) {
	h, ok := w.handlers[ethcommon.HexToAddress(addr)]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return h, nil
}

func (w *EVMVoter) RegisterProposalHandler(address string, handler ProposalHandler) {
	w.handlers[ethcommon.HexToAddress(address)] = handler
}
