package writer

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/ChainSafe/chainbridgev2/relayer"
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
	QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error)
	VoterAccountID() types.AccountID
}

type ProposalHandler func(msg *relayer.Message) (*SubstrateProposal, error)
type ProposalHandlers map[relayer.TransferType]ProposalHandler

type SubstrateWriter struct {
	client   Voter
	handlers ProposalHandlers
	chainID  uint8
}

func NewSubstrateWriter(chainID uint8, client Voter) *SubstrateWriter {
	return &SubstrateWriter{chainID: chainID, client: client}
}

func (w *SubstrateWriter) RegisterHandler(t relayer.TransferType, handler ProposalHandler) {
	w.handlers[t] = handler
}

func (w *SubstrateWriter) VoteProposal(m *relayer.Message) error {
	handler, ok := w.handlers[m.Type]
	if !ok {
		return errors.New(fmt.Sprintf("no corresponding substrate handler found for message type %s", m.Type))
	}
	prop, err := handler(m)

	if err != nil {
		return fmt.Errorf("failed to construct proposal (chain=%d, name=%s) Error: %w", m.Destination, w.chainID, err)
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
	var voteState struct {
		VotesFor     []types.AccountID
		VotesAgainst []types.AccountID
		Status       struct {
			IsActive   bool
			IsApproved bool
			IsRejected bool
		}
	}

	voteRes := &voteState
	srcId, err := types.EncodeToBytes(prop.SourceId)
	if err != nil {
		return false, "", err
	}
	propBz, err := prop.Encode()
	if err != nil {
		return false, "", err
	}
	exists, err := w.client.QueryStorage(BridgeStoragePrefix, "Votes", srcId, propBz, &voteRes)
	if err != nil {
		return false, "", err
	}

	if !exists {
		return true, "", nil
	} else if voteRes.Status.IsActive {
		if containsVote(voteRes.VotesFor, w.client.VoterAccountID()) ||
			containsVote(voteRes.VotesAgainst, w.client.VoterAccountID()) {
			return false, "already voted", nil
		} else {
			return true, "", nil
		}
	} else {
		return false, "proposal complete", nil
	}
}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
