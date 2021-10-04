// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package voter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
)

type ChainClient interface {
	LatestBlock() (*big.Int, error)
	RelayerAddress() common.Address
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
	ChainID(ctx context.Context) (*big.Int, error)
	calls.ClientDispatcher
}

type Proposer interface {
	Status(client ChainClient) (relayer.ProposalStatus, error)
	VotedBy(client ChainClient, by common.Address) (bool, error)
	Vote(client ChainClient, fabric calls.TxFabric) error
}

type MessageHandler interface {
	HandleMessage(m *relayer.Message) (Proposer, error)
}

type EVMVoter struct {
	stop   <-chan struct{}
	mh     MessageHandler
	client ChainClient
	fabric calls.TxFabric
}

func NewVoter(mh MessageHandler, client ChainClient, fabric calls.TxFabric) *EVMVoter {
	return &EVMVoter{
		mh:     mh,
		client: client,
		fabric: fabric,
	}
}

func (w *EVMVoter) VoteProposal(m *relayer.Message) error {
	prop, err := w.mh.HandleMessage(m)
	if err != nil {
		return err
	}
	ps, err := prop.Status(w.client)
	if err != nil {
		return fmt.Errorf("error getting proposal: %+v status %w", prop, err)
	}
	votedByTheRelayer, err := prop.VotedBy(w.client, w.client.RelayerAddress())
	if err != nil {
		return err
	}
	// if this relayer had not voted for proposal and proposal in Active status then we need to vote for
	// And that basically it o other options compared to previous contracts version
	if !votedByTheRelayer && ps == relayer.ProposalStatusActive {
		err = prop.Vote(w.client, w.fabric)
		if err != nil {
			return fmt.Errorf("Voting failed. Err: %w", err)
		}
	}
	return nil
}
