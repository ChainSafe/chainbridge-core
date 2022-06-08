// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/store"
	"github.com/rs/zerolog/log"
)

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int, msgChan chan *message.Message, errChan chan<- error)
}

type ProposalExecutor interface {
	Execute(message *message.Message) error
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener   EventListener
	writer     ProposalExecutor
	blockstore *store.BlockStore
	config     *chain.EVMConfig
}

func NewEVMChain(listener EventListener, writer ProposalExecutor, blockstore *store.BlockStore, config *chain.EVMConfig) *EVMChain {
	return &EVMChain{listener: listener, writer: writer, blockstore: blockstore, config: config}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *EVMChain) PollEvents(ctx context.Context, sysErr chan<- error, msgChan chan *message.Message) {
	log.Info().Msg("Polling Blocks...")

	startBlock, err := c.blockstore.GetStartBlock(
		*c.config.GeneralChainConfig.Id,
		c.config.StartBlock,
		c.config.GeneralChainConfig.LatestBlock,
		c.config.GeneralChainConfig.FreshStart,
	)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}

	go c.listener.ListenToEvents(ctx, startBlock, msgChan, sysErr)
}

func (c *EVMChain) Write(msg *message.Message) error {
	return c.writer.Execute(msg)
}

func (c *EVMChain) DomainID() uint8 {
	return *c.config.GeneralChainConfig.Id
}
