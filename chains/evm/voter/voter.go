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
	"github.com/ChainSafe/chainbridge-core/relayer"
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
	HandleMessage(m *relayer.Message) (*proposal.Proposal, error)
}

type EVMVoter struct {
	mh                   MessageHandler
	client               ChainClient
	fabric               calls.TxFabric
	gasPriceClient       calls.GasPricer
	pendingProposalVotes map[common.Hash]uint8
}

func NewVoter(mh MessageHandler, client ChainClient, fabric calls.TxFabric, gasPriceClient calls.GasPricer) (*EVMVoter, error) {
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

// VoteProposal checks if relayer already voted and is threshold
// satisfied and casts a vote if it isn't
func (v *EVMVoter) VoteProposal(m *relayer.Message) error {
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

	shouldVoteChn := make(chan bool)
	go v.shouldVoteForProposal(shouldVoteChn, prop, 0)

	shouldVote := <-shouldVoteChn
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

func (v *EVMVoter) shouldVoteForProposal(shouldVote chan bool, prop *proposal.Proposal, tries int) {
	propID := prop.GetID()
	defer delete(v.pendingProposalVotes, propID)

	Sleep(time.Duration(rand.Intn(shouldVoteCheckPeriod)) * time.Second)

	ps, err := calls.ProposalStatus(v.client, prop)
	if err != nil {
		log.Error().Err(err)
		shouldVote <- false
		return
	}

	if ps.Status == relayer.ProposalStatusExecuted || ps.Status == relayer.ProposalStatusCanceled {
		shouldVote <- false
		return
	}

	threshold, err := calls.GetThreshold(v.client, &prop.BridgeAddress)
	if err != nil {
		log.Error().Err(err)
		shouldVote <- false
		return
	}

	if ps.YesVotesTotal+v.pendingProposalVotes[propID] >= threshold && tries < maxShouldVoteChecks {
		// Wait until proposal status is finalized to prevent missing votes
		// in case of dropped txs
		tries++
		v.shouldVoteForProposal(shouldVote, prop, tries)
		return
	}

	shouldVote <- true
}

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

func (v *EVMVoter) increaseProposalVoteCount(hash common.Hash, propID common.Hash) {
	v.pendingProposalVotes[propID]++

	_, err := v.client.WaitAndReturnTxReceipt(hash)
	if err != nil {
		log.Error().Err(err)
	}

	v.pendingProposalVotes[propID]--
}
