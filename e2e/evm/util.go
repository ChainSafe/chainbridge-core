package evm

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/consts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/events"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var TestTimeout = time.Second * 600

type Client interface {
	LatestBlock() (*big.Int, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error)
}

func WaitUntilProposalExecuted(client Client, bridge common.Address) error {
	startBlock, _ := client.LatestBlock()

	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		Addresses: []common.Address{bridge},
		Topics: [][]common.Hash{
			{events.ProposalEventSig.GetTopic()},
		},
	}
	timeout := time.After(TestTimeout)
	ch := make(chan types.Log)
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return err
	}
	sub, err := client.SubscribeFilterLogs(context.Background(), query, ch)
	// if unable to subscribe check for the proposal execution every 5 sec
	if err != nil {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				endBlock, _ := client.LatestBlock()
				res, err := checkProposalExecuted(client, startBlock, endBlock, bridge, a)
				if err != nil {
					return err
				}
				if res {
					return nil
				}
				startBlock = endBlock
			case <-timeout:
				return errors.New("test timed out waiting for ProposalCreated event")
			}
		}
	}
	defer sub.Unsubscribe()

	for {
		select {
		case evt := <-ch:
			out, err := a.Unpack("ProposalEvent", evt.Data)
			if err != nil {
				return err
			}
			status := abi.ConvertType(out[2], new(uint8)).(*uint8)
			// Check status
			if IsExecuted(*status) {
				log.Info().Msgf("Got Proposal executed event status, continuing..., status: %v", *status)
				return nil
			} else {
				log.Info().Msgf("Got Proposal event status: %v", *status)
			}
		case err := <-sub.Err():
			if err != nil {
				return err
			}
		case <-timeout:
			return errors.New("test timed out waiting for ProposalCreated event")
		}
	}
}

type ProposalStatus int

const (
	Inactive ProposalStatus = iota
	Active
	Passed
	Executed
	Cancelled
)

func IsActive(status uint8) bool {
	return ProposalStatus(status) == Active
}

func IsFinalized(status uint8) bool {
	return ProposalStatus(status) == Passed
}

func IsExecuted(status uint8) bool {
	return ProposalStatus(status) == Executed
}

func checkProposalExecuted(client Client, startBlock, endBlock *big.Int, bridge common.Address, a abi.ABI) (bool, error) {
	logs, err := client.FetchEventLogs(context.TODO(), bridge, string(events.ProposalEventSig), startBlock, endBlock)
	if err != nil {
		return false, err
	}
	for _, evt := range logs {
		out, err := a.Unpack("ProposalEvent", evt.Data)
		if err != nil {
			return false, err
		}
		status := abi.ConvertType(out[2], new(uint8)).(*uint8)
		if IsExecuted(*status) {
			log.Info().Msgf("Got Proposal executed event status, continuing..., status: %v", *status)
			return true, nil
		} else {
			log.Info().Msgf("Got Proposal event status: %v", *status)
		}
	}
	return false, nil
}

func WaitUntilBridgeReady(client Client, bridge common.Address) error {
	startBlock, _ := client.LatestBlock()
	logs, err := client.FetchEventLogs(context.Background(), bridge, string(events.ThresholdChangedSig), big.NewInt(1), startBlock)
	if err != nil {
		return err
	}
	if len(logs) > 0 {
		return nil
	}

	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		Addresses: []common.Address{bridge},
		Topics: [][]common.Hash{
			{events.ThresholdChangedSig.GetTopic()},
		},
	}
	ch := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, ch)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case <-ch:
			return nil
		case err := <-sub.Err():
			if err != nil {
				return err
			}
		}
	}
}
