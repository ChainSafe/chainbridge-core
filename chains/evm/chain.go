// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/rs/zerolog/log"
)

type EventListener interface {
	ListenToEvents(startBlock *big.Int, domainID uint8, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan *message.Message
}

type ProposalVoter interface {
	VoteProposal(message *message.Message) error
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener EventListener // Rename
	writer   ProposalVoter
	domainID uint8
	kvdb     blockstore.KeyValueReaderWriter
	config   *chain.SharedEVMConfig
}

func NewEVMChain(dr EventListener, writer ProposalVoter, kvdb blockstore.KeyValueReaderWriter, domainID uint8, config *chain.SharedEVMConfig) *EVMChain {
	return &EVMChain{listener: dr, writer: writer, kvdb: kvdb, domainID: domainID, config: config}
}

// PollEvents is the goroutine that polling blocks and searching Deposit Events in them. Event then sent to eventsChan
func (c *EVMChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan *message.Message) {
	log.Info().Msg("Polling Blocks...")
	// Handler chain specific configs and flags
	block, err := blockstore.SetupBlockstore(&c.config.GeneralChainConfig, c.kvdb, c.config.StartBlock)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}
	ech := c.listener.ListenToEvents(block, c.domainID, c.kvdb, stop, sysErr)
	for {
		select {
		case <-stop:
			return
		case newEvent := <-ech:
			// Here we can place middlewares for custom logic?
			eventsChan <- newEvent
			continue
		}
	}
}

func (c *EVMChain) Write(msg *message.Message) error {
	return c.writer.VoteProposal(msg)
}

func (c *EVMChain) DomainID() uint8 {
	return c.domainID
}
