package evm

import (
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/rs/zerolog/log"
)

var ErrFatalPolling = errors.New("listener block polling failed")
var BlockRetryLimit = 5
var BlockRetryInterval = time.Second * 5
var BlockDelay = big.NewInt(10) //TODO: move to config

//type AddressType [32]byte

// DepositReader
type EventReader interface {
	//LatestBlock() (*big.Int, error)
	//MatchResourceIDToHandlerAddress(rID [32]byte) (ethcommon.Address, error)
	//MatchAddressWithHandlerFunc(addr ethcommon.Address) (Handler, error)
	//LogsForBlock(ctx context.Context, latestBlock *big.Int) ([]types.Log, error)
	//GetContractBackend() bind.ContractBackend
	GetDepositEventsForBlockRange(blockFrom, blockTo *big.Int) ([]relayer.XCMessager, error)
}

type LatestBlockGetter interface {
	LatestBlock() (*big.Int, error)
}

type EVMWriter interface {
	Write()
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener EventReader // Rename
	writer   EVMWriter
	chainID  uint8
	block    *big.Int
	bg       LatestBlockGetter
	bs       relayer.BlockStorer
}

func NewEVMChain(dr EventReader, writer EVMWriter, bs relayer.BlockStorer) *EVMChain {
	return &EVMChain{listener: dr, writer: writer, bs: bs}
}

// PollEvents is the gorutine that polling blocks and searching Deposit Events in them. Event then sent to eventsChan
func (c *EVMChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan relayer.XCMessager) {
	log.Info().Msg("Polling Blocks...")
	for {
		select {
		case <-stop:
			return
		default:
			latestBlock, err := c.bg.LatestBlock()
			if err != nil {
				log.Error().Err(err).Str("block", c.block.String()).Msg("Unable to get latest block")
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, c.block).Cmp(BlockDelay) == -1 {
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Parse out events
			events, err := c.listener.GetDepositEventsForBlockRange(c.block, c.block)
			if err != nil {
				log.Error().Str("block", c.block.String()).Err(err).Msg("Failed to get events for block")
				time.Sleep(BlockRetryInterval)
				continue
			}
			// TODO: FIX THIS
			for _, e := range events {
				eventsChan <- e
			}

			if c.block.Int64()%20 == 0 {
				// Logging process every 20 bocks to exclude spam
				log.Debug().Str("block", c.block.String()).Msg("Queried block for deposit events")
			}

			//Write to block store. Not a critical operation, no need to retry
			err = c.bs.StoreBlock(c.block, c.chainID)
			if err != nil {
				log.Error().Str("block", c.block.String()).Err(err).Msg("Failed to write latest block to blockstore")
			}

			// Goto next block
			c.block.Add(c.block, big.NewInt(1))
		}
	}
}

// Write function pass XCMessager to underlying chain writer
func (c *EVMChain) Write(relayer.XCMessager) {
	c.writer.Write() // TODO
}

func (c *EVMChain) ChainID() uint8 {
	return c.chainID
}
