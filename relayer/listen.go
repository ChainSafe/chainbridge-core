package relayer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type Handler func(sourceID, destID uint8, nonce uint64, handlerContractAddress common.Address, backend bind.ContractBackend) (*XCMessage, error)
type Handlers map[ethcommon.Address]Handler

type IListener interface {
	LatestBlock() (*big.Int, error)
	GetBridgeAddress() ethcommon.Address
	StoreCurrentBlock(*big.Int) error
	MatchResourceIDToHandlerAddress(rID [32]byte) (ethcommon.Address, error)
	MatchAddressWithHandler(addr ethcommon.Address) (Handler, error)
	LogsForBlock(ctx context.Context, latestBlock *big.Int) ([]types.Log, error)
	GetContractBackend() bind.ContractBackend
	GetChainID() uint8
}

var ErrFatalPolling = errors.New("listener block polling failed")
var BlockRetryLimit = 5
var BlockRetryInterval = time.Second * 5
var BlockDelay = big.NewInt(1) //TODO: move to config

func PollBlocks(l IListener, stop <-chan struct{}, sysErr chan<- error) {
	log.Info().Msg("Polling Blocks...")
	var currentBlock = big.NewInt(0)
	var retry = BlockRetryLimit
	for {
		select {
		case <-stop:
			return
		default:
			// No more retries, goto next block
			if retry == 0 {
				log.Error().Msg("Polling failed, retries exceeded")
				sysErr <- ErrFatalPolling
				return
			}
			latestBlock, err := l.LatestBlock()
			if err != nil {
				log.Error().Err(err).Str("block", currentBlock.String()).Msg("Unable to get latest block")
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, currentBlock).Cmp(BlockDelay) == -1 {
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Parse out events
			err = getDepositEventsForBlock(l, currentBlock)
			if err != nil {
				log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to get events for block")
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}
			if currentBlock.Int64()%20 == 0 {
				// Logging process every 20 bocks to exclude spam
				log.Debug().Str("block", currentBlock.String()).Msg("Queried block for deposit events")
			}

			//Write to block store. Not a critical operation, no need to retry
			err = l.StoreCurrentBlock(currentBlock)
			if err != nil {
				log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
			}

			// Goto next block and reset retry counter
			currentBlock.Add(currentBlock, big.NewInt(1))
			retry = BlockRetryLimit
		}
	}
}

const (
	Deposit string = "Deposit(uint8,bytes32,uint64)"
)

func getDepositEventsForBlock(l IListener, latestBlock *big.Int) error {
	logs, err := l.LogsForBlock(context.Background(), latestBlock)
	if err != nil {
		return fmt.Errorf("unable to Filter Logs: %w", err)
	}
	if len(logs) == 0 {
		return nil
	}
	// read through the log events and handle their deposit event if handler is recognized
	for _, eventLog := range logs {
		destId := uint8(eventLog.Topics[1].Big().Uint64())
		rId := eventLog.Topics[2]
		nonce := eventLog.Topics[3].Big().Uint64()

		addr, err := l.MatchResourceIDToHandlerAddress(rId)
		if err != nil {
			return err
		}

		eventHandler, err := l.MatchAddressWithHandler(addr)
		if err != nil {
			return fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rId, err)
		}
		backend := l.GetContractBackend()
		m, err := eventHandler(l.GetChainID(), destId, nonce, addr, backend)
		log.Debug().Msgf("Resolved message %+v", m)
		// Here we should send message to dest writer. For this we need to have router instance, but it will require to pass it, so maybe we can return some channel and send events to channel
	}
	return nil
}
