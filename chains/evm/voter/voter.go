// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package voter

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

const (
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
	calls.ClientDispatcher
}

type MessageHandler interface {
	HandleMessage(m *message.Message) (*proposal.Proposal, error)
}

type EVMVoter struct {
	mh                   MessageHandler
	client               ChainClient
	fabric               calls.TxFabric
	gasPriceClient       calls.GasPricer
	pendingProposalVotes map[common.Hash]uint8
}

// NewVoterWithSubscription creates EVMVoter with pending proposal subscription
// that listens to pending voteProposal transactions and avoids wasting gas
// on sending votes for transactions that will fail.
// Currently, officialy supported only by Geth nodes.
func NewVoterWithSubscription(mh MessageHandler, client ChainClient, fabric calls.TxFabric, gasPriceClient calls.GasPricer) (*EVMVoter, error) {
	voter := &EVMVoter{
		mh:                   mh,
		client:               client,
		fabric:               fabric,
		gasPriceClient:       gasPriceClient,
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

// NewVoter creates EVMVoter without pending proposal subscription.
// It is a fallback for nodes that don't support pending transaction subscription
// and will vote on proposals that already satisfy threshold thus wasting gas.
func NewVoter(mh MessageHandler, client ChainClient, fabric calls.TxFabric, gasPriceClient calls.GasPricer) *EVMVoter {
	return &EVMVoter{
		mh:                   mh,
		client:               client,
		fabric:               fabric,
		gasPriceClient:       gasPriceClient,
		pendingProposalVotes: make(map[common.Hash]uint8),
	}
}

// VoteProposal checks if relayer already voted and is threshold
// satisfied and casts a vote if it isn't.
func (v *EVMVoter) VoteProposal(m *message.Message) error {
	prop, err := v.mh.HandleMessage(m)
	if err != nil {
		return err
	}

	votedByTheRelayer, err := calls.IsProposalVotedBy(v.client, v.client.RelayerAddress(), prop)
	if err != nil {
		return err
	}
	if votedByTheRelayer {
		return nil
	}

	shouldVote, err := v.shouldVoteForProposal(prop, 0)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if !shouldVote {
		log.Debug().Msgf("Proposal %+v already satisfies threshold", prop)
		return nil
	}

	hash, err := calls.VoteProposal(v.client, v.fabric, v.gasPriceClient, prop)
	if err != nil {
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

	ps, err := calls.ProposalStatus(v.client, prop)
	if err != nil {
		return false, err
	}

	if ps.Status == message.ProposalStatusExecuted || ps.Status == message.ProposalStatusCanceled {
		return false, nil
	}

	threshold, err := calls.GetThreshold(v.client, &prop.BridgeAddress)
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
