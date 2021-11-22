// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
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
	listener EventListener
	writer   ProposalVoter
	kvdb     blockstore.KeyValueReaderWriter
	config   *chain.EVMConfig
}

// SetupDefaultEVMChain sets up an EVMChain with all supported handlers configured
func SetupDefaultEVMChain(rawConfig map[string]interface{}, db blockstore.KeyValueReaderWriter) (*EVMChain, error) {
	config, err := chain.NewEVMConfig(rawConfig)
	if err != nil {
		return nil, err
	}

	client := evmclient.NewEVMClient()
	err = client.Configurate(config)
	if err != nil {
		return nil, err
	}

	eventHandler := listener.NewETHEventHandler(common.HexToAddress(config.Bridge), client)
	eventHandler.RegisterEventHandler(config.Erc20Handler, listener.Erc20EventHandler)
	eventHandler.RegisterEventHandler(config.Erc721Handler, listener.Erc721EventHandler)
	eventHandler.RegisterEventHandler(config.GenericHandler, listener.GenericEventHandler)
	evmListener := listener.NewEVMListener(client, eventHandler, common.HexToAddress(config.Bridge))

	mh := voter.NewEVMMessageHandler(client, common.HexToAddress(config.Bridge))
	mh.RegisterMessageHandler(config.Erc20Handler, voter.ERC20MessageHandler)
	mh.RegisterMessageHandler(config.Erc721Handler, voter.ERC721MessageHandler)
	mh.RegisterMessageHandler(config.GenericHandler, voter.GenericMessageHandler)
	evmVoter := voter.NewVoter(mh, client, evmtransaction.NewTransaction, evmgaspricer.NewLondonGasPriceClient(client, nil))

	return NewEVMChain(evmListener, evmVoter, db, config), nil
}

func NewEVMChain(listener EventListener, writer ProposalVoter, kvdb blockstore.KeyValueReaderWriter, config *chain.EVMConfig) *EVMChain {
	return &EVMChain{listener: listener, writer: writer, kvdb: kvdb, config: config}
}

// PollEvents is the goroutine that polls blocks and searches Deposit events in them.
// Events are then sent to eventsChan.
func (c *EVMChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan *message.Message) {
	log.Info().Msg("Polling Blocks...")

	startBlock, err := blockstore.GetStartBlock(
		c.kvdb,
		*c.config.GeneralChainConfig.Id,
		c.config.StartBlock,
		c.config.GeneralChainConfig.LatestBlock,
		c.config.GeneralChainConfig.FreshStart,
	)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}

	ech := c.listener.ListenToEvents(startBlock, *c.config.GeneralChainConfig.Id, c.kvdb, stop, sysErr)
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
	return *c.config.GeneralChainConfig.Id
}
