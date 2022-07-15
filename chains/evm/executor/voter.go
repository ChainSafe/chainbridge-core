// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package executor

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/consts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"

	"github.com/ChainSafe/sygma-core/chains/evm/executor/proposal"
	"github.com/ChainSafe/sygma-core/relayer/message"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

const (
	maxSimulateVoteChecks = 5
	maxShouldVoteChecks   = 40
	shouldVoteCheckPeriod = 15
)

var (
	Sleep = time.Sleep
)

type ChainClient interface {
	RelayerAddress() common.Address
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
	SubscribePendingTransactions(ctx context.Context, ch chan<- common.Hash) (*rpc.ClientSubscription, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *ethereumTypes.Transaction, isPending bool, err error)
	calls.ContractCallerDispatcher
}

type MessageHandler interface {
	HandleMessage(m *message.Message) (*proposal.Proposal, error)
}

type BridgeContract interface {
	IsProposalVotedBy(by common.Address, p *proposal.Proposal) (bool, error)
	VoteProposal(proposal *proposal.Proposal, opts transactor.TransactOptions) (*common.Hash, error)
	SimulateVoteProposal(proposal *proposal.Proposal) error
	ProposalStatus(p *proposal.Proposal) (message.ProposalStatus, error)
	GetThreshold() (uint8, error)
}

type EVMVoter struct {
	mh                   MessageHandler
	client               ChainClient
	bridgeContract       BridgeContract
	pendingProposalVotes map[common.Hash]uint8
}

// NewVoterWithSubscription creates an instance of EVMVoter that votes for
// proposals on chain.
//
// It is created with a pending proposal subscription that listens to
// pending voteProposal transactions and avoids wasting gas on sending votes
// for transactions that will fail.
// Currently, officially supported only by Geth nodes.
func NewVoterWithSubscription(mh MessageHandler, client ChainClient, bridgeContract BridgeContract) (*EVMVoter, error) {
	voter := &EVMVoter{
		mh:                   mh,
		client:               client,
		bridgeContract:       bridgeContract,
		pendingProposalVotes: make(map[common.Hash]uint8),
	}

	ch := make(chan common.Hash)
	_, err := client.SubscribePendingTransactions(context.TODO(), ch)
	if err != nil {
		return nil, err
	}
	go voter.trackProposalPendingVotes(ch)

	return voter, nil
}

// NewVoter creates an instance of EVMVoter that votes for proposal on chain.
//
// It is created without pending proposal subscription and is a fallback
// for nodes that don't support pending transaction subscription and will vote
// on proposals that already satisfy threshold.
func NewVoter(mh MessageHandler, client ChainClient, bridgeContract BridgeContract) *EVMVoter {
	return &EVMVoter{
		mh:                   mh,
		client:               client,
		bridgeContract:       bridgeContract,
		pendingProposalVotes: make(map[common.Hash]uint8),
	}
}

// Execute checks if relayer already voted and is threshold
// satisfied and casts a vote if it isn't.
func (v *EVMVoter) Execute(m *message.Message) error {
	prop, err := v.mh.HandleMessage(m)
	if err != nil {
		return err
	}

	votedByTheRelayer, err := v.bridgeContract.IsProposalVotedBy(v.client.RelayerAddress(), prop)
	if err != nil {
		log.Error().Err(err).Msgf("Fetching is proposal %v voted by relayer failed", prop)
		return err
	}
	if votedByTheRelayer {
		return nil
	}

	shouldVote, err := v.shouldVoteForProposal(prop, 0)
	if err != nil {
		log.Error().Err(err).Msgf("Should vote for proposal %v failed", prop)
		return err
	}

	if !shouldVote {
		log.Debug().Msgf("Proposal %+v already satisfies threshold", prop)
		return nil
	}
	err = v.repetitiveSimulateVote(prop, 0)
	if err != nil {
		log.Error().Err(err).Msgf("Simulating proposal %+v vote failed", prop)
		return err
	}

	hash, err := v.bridgeContract.VoteProposal(prop, transactor.TransactOptions{Priority: prop.Metadata.Priority})
	if err != nil {
		log.Error().Err(err).Msgf("voting for proposal %+v failed", prop)
		return fmt.Errorf("voting failed. Err: %w", err)
	}

	log.Debug().Str("hash", hash.String()).Uint64("nonce", prop.DepositNonce).Msgf("Voted")
	return nil
}

// shouldVoteForProposal checks if proposal already has threshold with pending
// proposal votes from other relayers.
// Only works properly in conjuction with NewVoterWithSubscription as without a subscription
// no pending txs would be received and pending vote count would be 0.
func (v *EVMVoter) shouldVoteForProposal(prop *proposal.Proposal, tries int) (bool, error) {
	propID := prop.GetID()
	defer delete(v.pendingProposalVotes, propID)

	// random delay to prevent all relayers checking for pending votes
	// at the same time and all of them sending another tx
	Sleep(time.Duration(rand.Intn(shouldVoteCheckPeriod)) * time.Second)

	ps, err := v.bridgeContract.ProposalStatus(prop)
	if err != nil {
		return false, err
	}

	if ps.Status == message.ProposalStatusExecuted || ps.Status == message.ProposalStatusCanceled {
		return false, nil
	}

	threshold, err := v.bridgeContract.GetThreshold()
	if err != nil {
		return false, err
	}

	if ps.YesVotesTotal+v.pendingProposalVotes[propID] >= threshold && tries < maxShouldVoteChecks {
		// Wait until proposal status is finalized to prevent missing votes
		// in case of dropped txs
		tries++
		return v.shouldVoteForProposal(prop, tries)
	}

	return true, nil
}

// repetitiveSimulateVote repeatedly tries(5 times) to simulate vore proposal call until it succeeds
func (v *EVMVoter) repetitiveSimulateVote(prop *proposal.Proposal, tries int) error {
	err := v.bridgeContract.SimulateVoteProposal(prop)
	if err != nil {
		if tries < maxSimulateVoteChecks {
			tries++
			return v.repetitiveSimulateVote(prop, tries)
		}
		return err
	} else {
		return nil
	}
}

// trackProposalPendingVotes tracks pending voteProposal txs from
// other relayers and increases count of pending votes in pendingProposalVotes map
// by proposal unique id.
func (v *EVMVoter) trackProposalPendingVotes(ch chan common.Hash) {
	for msg := range ch {
		txData, _, err := v.client.TransactionByHash(context.TODO(), msg)
		if err != nil {
			log.Error().Err(err)
			continue
		}

		a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
		if err != nil {
			log.Error().Err(err)
			continue
		}

		if len(txData.Data()) < 4 {
			continue
		}

		m, err := a.MethodById(txData.Data()[:4])
		if err != nil {
			continue
		}

		data, err := m.Inputs.UnpackValues(txData.Data()[4:])
		if err != nil {
			log.Error().Err(err)
			continue
		}

		if m.Name == "voteProposal" {
			source := data[0].(uint8)
			depositNonce := data[1].(uint64)
			prop := proposal.Proposal{
				Source:       source,
				DepositNonce: depositNonce,
			}

			go v.increaseProposalVoteCount(msg, prop.GetID())
		}
	}
}

// increaseProposalVoteCount increases pending proposal vote for target proposal
// and decreases it when transaction is mined.
func (v *EVMVoter) increaseProposalVoteCount(hash common.Hash, propID common.Hash) {
	v.pendingProposalVotes[propID]++

	_, err := v.client.WaitAndReturnTxReceipt(hash)
	if err != nil {
		log.Error().Err(err)
	}

	v.pendingProposalVotes[propID]--
}
