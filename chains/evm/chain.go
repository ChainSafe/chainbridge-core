package evm

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridgev2/blockstore"
	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/rs/zerolog/log"
)

type EventListener interface {
	ListenToEvents(startBlock *big.Int, chainID uint8, bridgeContractAddress string, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan relayer.XCMessager
}

type ProposalVoter interface {
	VoteProposal(message relayer.XCMessager, bridgeAddress string) error
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener              EventListener // Rename
	writer                ProposalVoter
	chainID               uint8
	kvdb                  blockstore.KeyValueReaderWriter
	bridgeContractAddress string
}

func NewEVMChain(dr EventListener, writer ProposalVoter, kvdb blockstore.KeyValueReaderWriter, bridgeContractAddress string) *EVMChain {
	return &EVMChain{listener: dr, writer: writer, kvdb: kvdb, bridgeContractAddress: bridgeContractAddress}
}

// PollEvents is the goroutine that polling blocks and searching Deposit Events in them. Event then sent to eventsChan
func (c *EVMChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan relayer.XCMessager) {
	log.Info().Msg("Polling Blocks...")
	// Handler chain specific configs and flags
	b, err := blockstore.GetLastStoredBlock(c.kvdb, c.chainID)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}
	ech := c.listener.ListenToEvents(b, c.chainID, c.bridgeContractAddress, c.kvdb, stop, sysErr)
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

// Write function pass XCMessager to underlying chain writer
func (c *EVMChain) Write(msg relayer.XCMessager) error {
	return c.writer.VoteProposal(msg, c.bridgeContractAddress)
}

func (c *EVMChain) ChainID() uint8 {
	return c.chainID
}
