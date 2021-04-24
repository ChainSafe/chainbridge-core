package listener

import (
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridgev2/blockstore"
	"github.com/ChainSafe/chainbridgev2/chains/substrate"
	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/rs/zerolog/log"
)

const (
	FungibleTransfer    string = "FungibleTransfer"
	NonFungibleTransfer string = "NonFungibleTransfer"
	GenericTransfer     string = "GenericTransfer"
)

var BlockRetryInterval = time.Second * 5

var ErrBlockNotReady = errors.New("required result to be 32 bytes, but got 0")

type SubstrateClienter interface {
	GetHeaderLatest() (*types.Header, error)
	GetBlockHash(blockNumber uint64) (types.Hash, error)
	GetBlockEvents(hash types.Hash) (interface{}, error)
	UpdateMetatdata() error
}

type EventHandler func(interface{}) (*relayer.Message, error)

type SubstrateListener struct {
	conn          SubstrateClienter
	Subscriptions map[string]EventHandler
}

func (l *SubstrateListener) ListenToEvents(startBlock *big.Int, chainID uint8, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan *relayer.Message {
	ch := make(chan *relayer.Message)
	go func() {
		for {
			select {
			case <-stopChn:
				return
			default:
				if startBlock.Int64()%20 == 0 {
					// Logging process every 20 bocks to exclude spam
					log.Debug().Str("block", startBlock.String()).Uint8("chainID", chainID).Msg("Queried block for deposit events")
				}

				// retrieves the header of the latest block
				finalizedHeader, err := l.conn.GetHeaderLatest()
				if err != nil {
					log.Error().Err(err).Msg("Failed to fetch finalized header")
					time.Sleep(BlockRetryInterval)
					continue
				}
				if startBlock.Cmp(big.NewInt(0).SetUint64(uint64(finalizedHeader.Number))) == 1 {
					log.Error().Err(err).Msg("Failed to fetch finalized header")
					time.Sleep(BlockRetryInterval)
					continue
				}
				hash, err := l.conn.GetBlockHash(startBlock.Uint64())
				if err != nil && err.Error() == ErrBlockNotReady.Error() {
					time.Sleep(BlockRetryInterval)
					continue
				} else if err != nil {
					log.Error().Err(err).Str("block", startBlock.String()).Msg("Failed to query latest block")
					time.Sleep(BlockRetryInterval)
					continue
				}

				evts, err := l.conn.GetBlockEvents(hash)
				if err != nil {
					log.Error().Err(err).Msg("Failed to process events in block")
					continue
				}
				e, ok := evts.(*substrate.Events)
				if !ok {
					log.Error().Msg("Error decoding events")
				}
				msg, err := l.handleEvents(e)
				if err != nil {
					log.Error().Err(err).Msg("Error handling substrate events")
				}
				for _, m := range msg {
					ch <- m
				}
				err = blockstore.StoreBlock(kvrw, startBlock, chainID)
				if err != nil {
					log.Error().Str("block", startBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
				}
				startBlock.And(startBlock, big.NewInt(1))
			}
		}
	}()
	return ch
}

// handleEvents calls the associated handler for all registered event types
func (l *SubstrateListener) handleEvents(evts *substrate.Events) ([]*relayer.Message, error) {
	msgs := make([]*relayer.Message, 0)
	if l.Subscriptions[FungibleTransfer] != nil {
		for _, evt := range evts.ChainBridge_FungibleTransfer {
			m, err := l.Subscriptions[FungibleTransfer](evt)
			if err != nil {
				return nil, err
			}
			msgs = append(msgs, m)
		}
	}
	if l.Subscriptions[NonFungibleTransfer] != nil {
		for _, evt := range evts.ChainBridge_NonFungibleTransfer {
			m, err := l.Subscriptions[NonFungibleTransfer](evt)
			if err != nil {
				return nil, err
			}
			msgs = append(msgs, m)

		}
	}
	if l.Subscriptions[GenericTransfer] != nil {
		for _, evt := range evts.ChainBridge_GenericTransfer {
			m, err := l.Subscriptions[GenericTransfer](evt)
			if err != nil {
				return nil, err
			}
			msgs = append(msgs, m)
		}
	}
	if len(evts.System_CodeUpdated) > 0 {
		err := l.conn.UpdateMetatdata()
		if err != nil {
			log.Error().Err(err).Msg("Unable to update Metadata")
			return nil, err
		}
	}
	return msgs, nil
}
