package writer

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/substrate"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	internalTypes "github.com/ChainSafe/chainbridge-core/types"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/rs/zerolog/log"
)

const BridgePalletName = "ChainBridge"
const BridgeStoragePrefix = "ChainBridge"

var BlockRetryInterval = time.Second * 5
var BlockRetryLimit = 5
var AcknowledgeProposal = BridgePalletName + ".acknowledge_proposal"

type Voter interface {
	SubmitTx(method string, args ...interface{}) error
	GetVoterAccountID() types.AccountID
	GetMetadata() (meta types.Metadata)
	ResolveResourceId(resourceId internalTypes.ResourceID) (string, error)
	// TODO: Vote state should be higher abstraction
	GetProposalStatus(sourceID, proposalBytes []byte) (bool, *substrate.VoteState, error)
}

type ProposalHandler func(msg *message.Message) []interface{}
type ProposalHandlers map[message.TransferType]ProposalHandler

type SubstrateWriter struct {
	client   Voter
	handlers ProposalHandlers
	domainID uint8
}

func NewSubstrateWriter(domainID uint8, client Voter) *SubstrateWriter {
	return &SubstrateWriter{domainID: domainID, client: client}
}

func (w *SubstrateWriter) RegisterHandler(t message.TransferType, handler ProposalHandler) {
	if w.handlers == nil {
		w.handlers = make(map[message.TransferType]ProposalHandler)
	}
	w.handlers[t] = handler
}

func (w *SubstrateWriter) VoteProposal(m *message.Message) error {
	handler, ok := w.handlers[m.Type]
	if !ok {
		return fmt.Errorf("no corresponding substrate handler found for message type %s", m.Type)
	}
	prop, err := w.createProposal(m.Source, m.DepositNonce, m.ResourceId, handler(m)...)
	if err != nil {
		return fmt.Errorf("failed to construct proposal (chain=%d, name=%v) Error: %w", m.Destination, w.domainID, err)
	}

	for i := 0; i < BlockRetryLimit; i++ {
		// Ensure we only submit a vote if the proposal hasn't completed
		valid, reason, err := w.proposalValid(prop)
		if err != nil {
			time.Sleep(BlockRetryInterval)
			continue
		}

		// If active submit call, otherwise skip it. Retry on failure.
		if valid {
			err = w.client.SubmitTx(AcknowledgeProposal, prop.DepositNonce, prop.SourceId, prop.ResourceId, prop.Call)
			if err != nil {
				log.Error().Err(err).Msg("Failed to execute extrinsic")
				time.Sleep(BlockRetryInterval)
				continue
			}
			return nil
		} else {
			log.Info().Str("reason", reason).Uint64("nonce", uint64(prop.DepositNonce)).Uint8("source", uint8(prop.SourceId)).Str("resource", types.HexEncodeToString(prop.ResourceId[:])).Msg("Ignoring proposal")
			return nil
		}
	}
	return nil
}

func (w *SubstrateWriter) proposalValid(prop *SubstrateProposal) (bool, string, error) {
	srcId, err := types.EncodeToBytes(prop.SourceId)
	if err != nil {
		return false, "", err
	}
	propBz, err := prop.Encode()
	if err != nil {
		return false, "", err
	}
	exists, voteRes, err := w.client.GetProposalStatus(srcId, propBz)
	if err != nil {
		return false, "", err
	}
	if !exists {
		return true, "", nil
	} else if voteRes.Status.IsActive {
		if containsVote(voteRes.VotesFor, w.client.GetVoterAccountID()) ||
			containsVote(voteRes.VotesAgainst, w.client.GetVoterAccountID()) {
			return false, "already voted", nil
		} else {
			return true, "", nil
		}
	} else {
		return false, "proposal complete", nil
	}
}

func (w *SubstrateWriter) createProposal(sourceChain uint8, depositNonce uint64, resourceId internalTypes.ResourceID, args ...interface{}) (*SubstrateProposal, error) {
	meta := w.client.GetMetadata()
	method, err := w.client.ResolveResourceId(resourceId)
	if err != nil {
		return nil, err
	}
	call, err := types.NewCall(
		&meta,
		method,
		args...,
	)
	if err != nil {
		return nil, err
	}
	// TODO: Is not these should be always enabled?
	//if w.extendCall {
	//	eRID, err := types.EncodeToBytes(resourceId)
	//	if err != nil {
	//		return nil, err
	//	}
	//	call.Args = append(call.Args, eRID...)
	//}
	return &SubstrateProposal{
		DepositNonce: types.U64(depositNonce),
		Call:         call,
		SourceId:     types.U8(sourceChain),
		ResourceId:   types.NewBytes32(resourceId),
		Method:       method,
	}, nil
}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
