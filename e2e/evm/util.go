package evm

import (
	"context"
	"fmt"
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

func WaitForProposalExecuted(client TestClient, bridge common.Address) {
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
		panic(err)
	}
	defer sub.Unsubscribe()

	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		panic(err)
	}
	timeout := time.After(TestTimeout)
	for {
		select {
		case evt := <-ch:
			fmt.Println("nesjedosloooooo\nnn")
			out, err := a.Unpack("ProposalEvent", evt.Data)
			if err != nil {
				panic(err)
			}
			status := abi.ConvertType(out[2], new(uint8)).(*uint8)
			// Check status
			if util.IsExecuted(*status) {
				log.Info().Msgf("Got Proposal executed event status, continuing..., status: %v", *status)
				return
			} else {
				log.Info().Msgf("Got Proposal event status: %v", *status)
			}
		case err := <-sub.Err():
			if err != nil {
				panic(err)
			}
		case <-timeout:
			panic("Test timed out waiting for ProposalCreated event")
		}
	}
}
