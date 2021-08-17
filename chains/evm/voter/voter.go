// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package voter

import (
	"context"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var BlockRetryInterval = time.Second * 5

type ChainClient interface {
	LatestBlock() (*big.Int, error)
	SignAndSendTransaction(ctx context.Context, tx evmclient.CommonTransaction) (common.Hash, error)
	RelayerAddress() common.Address
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
	LockNonce()
	UnlockNonce()
	UnsafeNonce() (*big.Int, error)
	UnsafeIncreaseNonce() error
	BaseFee() (*big.Int, error)
	EstimateGasLondon(ctx context.Context, baseFee *big.Int) (*big.Int, *big.Int, error)
	GasPrice() (*big.Int, error)
	ChainID(ctx context.Context) (*big.Int, error)
}

type Proposer interface {
	Status(client ChainClient) (relayer.ProposalStatus, error)
	VotedBy(client ChainClient, by common.Address) (bool, error)
	Execute(client ChainClient, fabric TxFabric) error
	Vote(client ChainClient, fabric TxFabric) error
}

type MessageHandler interface {
	HandleMessage(m *relayer.Message) (Proposer, error)
}

type EVMVoter struct {
	stop   <-chan struct{}
	mh     MessageHandler
	client ChainClient
	fabric TxFabric
}

func NewVoter(mh MessageHandler, client ChainClient, fabric TxFabric) *EVMVoter {
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
		log.Error().Err(err).Msgf("error getting proposal status %+v", prop)
	}
	log.Debug().Msgf("Proposal status: %v", ps)

	votedByCurrentExecutor, err := prop.VotedBy(w.client, w.client.RelayerAddress())
	if err != nil {
		return err
	}
	log.Debug().Msgf("Voted by current executor: %v", votedByCurrentExecutor)

	if votedByCurrentExecutor || ps == relayer.ProposalStatusPassed || ps == relayer.ProposalStatusCanceled || ps == relayer.ProposalStatusExecuted {
		if ps == relayer.ProposalStatusPassed {
			// We should not vote for this proposal but it is ready to be executed
			err = prop.Execute(w.client, w.fabric)
			if err != nil {
				log.Error().Err(err).Msgf("Executing failed")
				return err
			}
			return nil
		} else {
			log.Debug().Bool("voted", votedByCurrentExecutor).Str("voter", w.client.RelayerAddress().String()).Msgf("proposal status %s", relayer.StatusMap[ps])
			return nil
		}
	}
	err = prop.Vote(w.client, w.fabric)
	if err != nil {
		log.Error().Err(err).Msgf("Voting failed")
		return err
	}
	// Checking every 5 seconds does proposal is ready to be executed
	// TODO: somehow update infinity loop to break after some period of time
	for {
		select {
		case <-time.After(BlockRetryInterval):
			ps, err := prop.Status(w.client)
			if err != nil {
				log.Error().Err(err).Msgf("error getting proposal status %+v", prop)
				return err
			}
			if ps == relayer.ProposalStatusPassed {
				err = prop.Execute(w.client, w.fabric)
				if err != nil {
					log.Error().Err(err).Msgf("Executing failed")
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

type DynamicGasPricer struct {
	client ChainClient
}

func NewDynamicGasPricer(client ChainClient) evmclient.GasPricer {
	return &DynamicGasPricer{client: client}
}

func (gasPricer *DynamicGasPricer) GasPrice() ([]*big.Int, error) {
	baseFee, err := gasPricer.client.BaseFee()
	if err != nil {
		return nil, err
	}

	var gasPrices []*big.Int
	if baseFee != nil {
		gasTipCap, gasFeeCap, err := gasPricer.client.EstimateGasLondon(context.TODO(), baseFee)
		if err != nil {
			return nil, err
		}
		gasPrices[0] = gasTipCap
		gasPrices[1] = gasFeeCap
	} else {
		gp, err := gasPricer.client.GasPrice()
		if err != nil {
			return nil, err
		}
		gasPrices[0] = gp
	}
	return gasPrices, nil
}
