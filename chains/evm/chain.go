package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

var ErrFatalPolling = errors.New("listener block polling failed")
var BlockRetryLimit = 5
var BlockRetryInterval = time.Second * 5
var BlockDelay = big.NewInt(10) //TODO: move to config

type Handler func(sourceID, destID uint8, nonce uint64, handlerContractAddress ethcommon.Address, backend bind.ContractCaller) (relayer.XCMessager, error)
type Handlers map[ethcommon.Address]Handler

type EVMListener interface {
	LatestBlock() (*big.Int, error)
	MatchResourceIDToHandlerAddress(rID [32]byte) (ethcommon.Address, error)
	MatchAddressWithHandlerFunc(addr ethcommon.Address) (Handler, error)
	LogsForBlock(ctx context.Context, latestBlock *big.Int) ([]types.Log, error)
	GetContractBackend() bind.ContractBackend
	GetChainID() uint8
}

type EVMWriter interface {
	Write()
}

type EVMChain struct {
	listener EVMListener
	writer   EVMWriter
}

func NewEVMChain(listener EVMListener) *EVMChain {
	return &EVMChain{}
}

func (c *EVMChain) PollEvents(bs relayer.BlockStorer, stop <-chan struct{}, sysErr chan<- error, eventsChan chan relayer.XCMessager) {
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
			latestBlock, err := c.listener.LatestBlock()
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
			err = getDepositEventsForBlock(c.listener, currentBlock, eventsChan)
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
			err = bs.StoreBlock(currentBlock, c.listener.GetChainID())
			if err != nil {
				log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
			}

			// Goto next block and reset retry counter
			currentBlock.Add(currentBlock, big.NewInt(1))
			retry = BlockRetryLimit
		}
	}
}

func getDepositEventsForBlock(l EVMListener, latestBlock *big.Int, eventsChan chan relayer.XCMessager) error {
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

		eventHandler, err := l.MatchAddressWithHandlerFunc(addr)
		if err != nil {
			return fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rId, err)
		}

		backend := l.GetContractBackend()
		m, err := eventHandler(l.GetChainID(), destId, nonce, addr, backend)
		log.Debug().Msgf("Resolved message %+v", m)

		// TODO: if noone to receive this will blocks forever
		eventsChan <- m
	}
	return nil
}

func (c *EVMChain) Write(relayer.XCMessager) {
	c.writer.Write() // TODO
}

func (c *EVMChain) ChainID() uint8 {
	return 0
}
