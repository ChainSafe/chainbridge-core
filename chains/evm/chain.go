package evm

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridgev2/blockstore"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/rs/zerolog/log"
)

type EventListener interface {
	ListenToEvents(startBlock *big.Int, stop <-chan struct{}, sysErr chan<- error) <-chan relayer.XCMessager
}

type LatestBlockGetter interface {
	LatestBlock() (*big.Int, error)
}

type EVMWriter interface {
	Write()
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener EventListener // Rename
	writer   EVMWriter
	chainID  uint8
	block    *big.Int
	bg       LatestBlockGetter
	kvdb     blockstore.KeyValueReaderWriter
}

func NewEVMChain(dr EventListener, writer EVMWriter, kvdb blockstore.KeyValueReaderWriter, bg LatestBlockGetter) *EVMChain {
	return &EVMChain{listener: dr, writer: writer, kvdb: kvdb, bg: bg}
}

// PollEvents is the goroutine that polling blocks and searching Deposit Events in them. Event then sent to eventsChan
func (c *EVMChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan relayer.XCMessager) {
	log.Info().Msg("Polling Blocks...")
	b, err := blockstore.GetLastStoredBlock(c.kvdb, c.chainID)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}
	c.block = b
	ech := c.listener.ListenToEvents(b, stop, sysErr)
	for {
		select {
		case <-stop:
			return
		case newEvent := <-ech:
			// Here we can place middlewares for custom logic
			eventsChan <- newEvent
			// TODO: We can store blocks to DB inside listener or make lestiener send something to channel each block to save it.
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
