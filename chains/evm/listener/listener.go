// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/ChainSafe/chainbridge-core/observability"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, startBlock *big.Int, endBlock *big.Int, msgChan chan []*message.Message) error
}

type ChainClient interface {
	LatestBlock() (*big.Int, error)
}

type BlockDeltaMeter interface {
	TrackBlockDelta(domainID uint8, head *big.Int, current *big.Int)
}

type EVMListener struct {
	client        ChainClient
	eventHandlers []EventHandler
	metrics       BlockDeltaMeter

	domainID           uint8
	blockstore         *store.BlockStore
	blockRetryInterval time.Duration
	blockConfirmations *big.Int
	blockInterval      *big.Int

	log zerolog.Logger
}

// NewEVMListener creates an EVMListener that listens to deposit events on chain
// and calls event handler when one occurs
func NewEVMListener(
	client ChainClient,
	eventHandlers []EventHandler,
	blockstore *store.BlockStore,
	metrics BlockDeltaMeter,
	domainID uint8,
	blockRetryInterval time.Duration,
	blockConfirmations *big.Int,
	blockInterval *big.Int) *EVMListener {
	logger := log.With().Uint8("domainID", domainID).Logger()
	return &EVMListener{
		log:                logger,
		client:             client,
		metrics:            metrics,
		eventHandlers:      eventHandlers,
		blockstore:         blockstore,
		domainID:           domainID,
		blockRetryInterval: blockRetryInterval,
		blockConfirmations: blockConfirmations,
		blockInterval:      blockInterval,
	}
}

// ListenToEvents goes block by block of a network and executes event handlers that are
// configured for the listener.
func (l *EVMListener) ListenToEvents(ctx context.Context, startBlock *big.Int, msgChan chan []*message.Message, errChn chan<- error) {
	endBlock := big.NewInt(0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ctxWithSpan, span, logger := observability.CreateSpanAndLoggerFromContext(ctx, "relayer-core", "relayer.core.EVMListener.ListenToEvents")
			head, err := l.client.LatestBlock()
			if err != nil {
				time.Sleep(l.blockRetryInterval)
				_ = observability.LogAndRecordError(&logger, span, err, "Unable to get latest block")
				span.End()
				continue
			}
			if startBlock == nil {
				startBlock = big.NewInt(head.Int64())
			}
			endBlock.Add(startBlock, l.blockInterval)

			// Sleep if the difference is less than needed block confirmations; (latest - current) < BlockDelay
			if new(big.Int).Sub(head, endBlock).Cmp(l.blockConfirmations) == -1 {
				time.Sleep(l.blockRetryInterval)
				observability.LogAndEvent(logger.Debug(), span, fmt.Sprintf("Block difference is less then %s", l.blockConfirmations))
				span.End()
				continue
			}

			l.metrics.TrackBlockDelta(l.domainID, head, endBlock)

			observability.LogAndEvent(logger.Debug(), span, fmt.Sprintf("Fetching evm events for block range %s-%s", startBlock, endBlock), attribute.String("startBlock", startBlock.String()), attribute.String("endBlock", endBlock.String()))

			for _, handler := range l.eventHandlers {
				err := handler.HandleEvent(ctxWithSpan, startBlock, new(big.Int).Sub(endBlock, big.NewInt(1)), msgChan)
				if err != nil {
					_ = observability.LogAndRecordError(&logger, span, err, "Unable to handle event")
					continue
				}
			}

			//Write to block store. Not a critical operation, no need to retry
			err = l.blockstore.StoreBlock(endBlock, l.domainID)
			if err != nil {
				_ = observability.LogAndRecordError(&logger, span, err, "Failed to write latest block to blockstore")
			}
			startBlock.Add(startBlock, l.blockInterval)
			span.End()
		}
	}
}
