// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChainSafe/sygma-core/store"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/rs/zerolog/log"
)

type EventListener interface {
	ListenToEvents(ctx context.Context, startBlock *big.Int, errChan chan<- error)
}

type ProposalExecutor interface {
	Execute(message *types.Message) error
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener   EventListener
	writer     ProposalExecutor
	blockstore *store.BlockStore

	domainID    uint8
	startBlock  *big.Int
	freshStart  bool
	latestBlock bool
}

func NewEVMChain(listener EventListener, writer ProposalExecutor, blockstore *store.BlockStore, domainID uint8, startBlock *big.Int, latestBlock bool, freshStart bool) *EVMChain {
	return &EVMChain{
		listener:    listener,
		writer:      writer,
		blockstore:  blockstore,
		domainID:    domainID,
		startBlock:  startBlock,
		latestBlock: latestBlock,
		freshStart:  freshStart,
	}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *EVMChain) PollEvents(ctx context.Context, sysErr chan<- error) {
	log.Info().Msg("Polling Blocks...")

	startBlock, err := c.blockstore.GetStartBlock(
		c.domainID,
		c.startBlock,
		c.latestBlock,
		c.freshStart,
	)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}

	go c.listener.ListenToEvents(ctx, startBlock, sysErr)
}

func (c *EVMChain) Write(msg []*types.Message) error {
	for _, msg := range msg {
		go func(msg *types.Message) {
			err := c.writer.Execute(msg)
			if err != nil {
				log.Err(err).Msgf("Failed writing message %v", msg)
			}
		}(msg)
	}

	return nil
}

func (c *EVMChain) DomainID() uint8 {
	return c.domainID
}
