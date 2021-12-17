package evm

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

var TestTimeout = time.Second * 600

func WaitForProposalExecuted(client TestClient, bridge common.Address) error {
	startBlock, _ := client.LatestBlock()

	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		Addresses: []common.Address{bridge},
		Topics: [][]common.Hash{
			{util.ProposalEvent.GetTopic()},
		},
	}
	ch := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, ch)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return err
	}
	timeout := time.After(TestTimeout)
	for {
		select {
		case evt := <-ch:
			out, err := a.Unpack("ProposalEvent", evt.Data)
			if err != nil {
				return err
			}
			status := abi.ConvertType(out[2], new(uint8)).(*uint8)
			// Check status
			if util.IsExecuted(*status) {
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
