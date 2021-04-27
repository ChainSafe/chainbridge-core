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

type TransferType string

const (
	FungibleTransfer    TransferType = "FungibleTransfer"
	NonFungibleTransfer TransferType = "NonFungibleTransfer"
	GenericTransfer     TransferType = "GenericTransfer"
)

var BlockRetryInterval = time.Second * 5

var ErrBlockNotReady = errors.New("required result to be 32 bytes, but got 0")

type SubstrateClienter interface {
	GetHeaderLatest() (*types.Header, error)
	GetBlockHash(blockNumber uint64) (types.Hash, error)
	GetBlockEvents(hash types.Hash, target interface{}) error
	UpdateMetatdata() error
}

type EventHandler func(interface{}) (*relayer.Message, error)

func NewSubstrateListener(client SubstrateClienter) *SubstrateListener {
	return &SubstrateListener{
		client: client,
	}
}

type SubstrateListener struct {
	client        SubstrateClienter
	subscriptions map[TransferType]EventHandler
}

func (l *SubstrateListener) RegisterHandler(tt TransferType, handler EventHandler) {
	l.subscriptions[tt] = handler
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
				finalizedHeader, err := l.client.GetHeaderLatest()
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
				hash, err := l.client.GetBlockHash(startBlock.Uint64())
				if err != nil && err.Error() == ErrBlockNotReady.Error() {
					time.Sleep(BlockRetryInterval)
					continue
				} else if err != nil {
					log.Error().Err(err).Str("block", startBlock.String()).Msg("Failed to query latest block")
					time.Sleep(BlockRetryInterval)
					continue
				}
				evts := &substrate.Events{}
				err = l.client.GetBlockEvents(evts)
				if err != nil {
					log.Error().Err(err).Msg("Failed to process events in block")
					continue
				}
				//e, ok := evts.(*substrate.Events)
				//if !ok {
				//	log.Error().Msg("Error decoding events")
				//}
				msg, err := l.handleEvents(evts)
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
	if l.subscriptions[FungibleTransfer] != nil {
		for _, evt := range evts.ChainBridge_FungibleTransfer {
			m, err := l.subscriptions[FungibleTransfer](evt)
			if err != nil {
				return nil, err
			}
			msgs = append(msgs, m)
		}
	}
	if l.subscriptions[NonFungibleTransfer] != nil {
		for _, evt := range evts.ChainBridge_NonFungibleTransfer {
			m, err := l.subscriptions[NonFungibleTransfer](evt)
			if err != nil {
				return nil, err
			}
			msgs = append(msgs, m)

		}
	}
	if l.subscriptions[GenericTransfer] != nil {
		for _, evt := range evts.ChainBridge_GenericTransfer {
			m, err := l.subscriptions[GenericTransfer](evt)
			if err != nil {
				return nil, err
			}
			msgs = append(msgs, m)
		}
	}
	if len(evts.System_CodeUpdated) > 0 {
		err := l.client.UpdateMetatdata()
		if err != nil {
			log.Error().Err(err).Msg("Unable to update Metadata")
			return nil, err
		}
	}
	return msgs, nil
}
