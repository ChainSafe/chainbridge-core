// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/store"

	"github.com/rs/zerolog/log"
)

type EventHandler interface {
	HandleEvent(block *big.Int, msgChan chan *message.Message) error
}

type ChainClient interface {
	LatestBlock() (*big.Int, error)
}

type EVMListener struct {
	client        ChainClient
	eventHandlers []EventHandler

	domainID           uint8
	blockstore         *store.BlockStore
	blockRetryInterval time.Duration
	blockConfirmations *big.Int
}

// NewEVMListener creates an EVMListener that listens to deposit events on chain
// and calls event handler when one occurs
func NewEVMListener(client ChainClient, eventHandlers []EventHandler, blockstore *store.BlockStore, config *chain.EVMConfig) *EVMListener {
	return &EVMListener{
		client:             client,
		eventHandlers:      eventHandlers,
		blockstore:         blockstore,
		domainID:           *config.GeneralChainConfig.Id,
		blockRetryInterval: config.BlockRetryInterval,
		blockConfirmations: config.BlockConfirmations,
	}
}

// ListenToEvents goes block by block of a network and executes event handlers that are
// configured for the listener.
func (l *EVMListener) ListenToEvents(ctx context.Context, block *big.Int, msgChan chan *message.Message, errChn chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			head, err := l.client.LatestBlock()
			if err != nil {
				log.Error().Err(err).Msg("Unable to get latest block")
				time.Sleep(l.blockRetryInterval)
				continue
			}
			if block == nil {
				block = head
			}
			// Sleep if the difference is less than needed block confirmations; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(head, block).Cmp(l.blockConfirmations) == -1 {
				time.Sleep(l.blockRetryInterval)
				continue
			}

			for _, handler := range l.eventHandlers {
				err := handler.HandleEvent(block, msgChan)
				if err != nil {
					log.Error().Err(err).Str("DomainID", string(l.domainID)).Msgf("Unable to handle events")
					continue
				}
			}

			//Write to block store. Not a critical operation, no need to retry
			err = l.blockstore.StoreBlock(block, l.domainID)
			if err != nil {
				log.Error().Str("block", block.String()).Err(err).Msg("Failed to write latest block to blockstore")
			}
			block.Add(block, big.NewInt(1))
		}
	}
}
