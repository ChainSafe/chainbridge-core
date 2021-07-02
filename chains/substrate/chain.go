package substrate

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/rs/zerolog/log"
)

type ProposalVoter interface {
	VoteProposal(message *relayer.Message) error
}

type EventListener interface {
	ListenToEvents(startBlock *big.Int, chainID uint8, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan *relayer.Message
}

type SubstrateChain struct {
	chainID  uint8
	stop     chan<- struct{}
	listener EventListener
	writer   ProposalVoter
	kvdb     blockstore.KeyValueReaderWriter
	config   *config.SharedSubstrateConfig
}

func NewSubstrateChain(listener EventListener, writer ProposalVoter, kvdb blockstore.KeyValueReaderWriter, chainID uint8, config *config.SharedSubstrateConfig) *SubstrateChain {
	return &SubstrateChain{
		listener: listener,
		writer:   writer,
		kvdb:     kvdb,
		chainID:  chainID,
		config:   config,
	}
}

func (c *SubstrateChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan *relayer.Message) {
	log.Info().Msg("Polling Blocks...")
	// Handler chain specific configs and flags
	//b, err := blockstore.GetLastStoredBlock(c.kvdb, c.chainID)
	block, err := blockstore.SetupBlockstore(&c.config.GeneralChainConfig, c.kvdb, c.config.StartBlock)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}
	ech := c.listener.ListenToEvents(block, c.chainID, c.kvdb, stop, sysErr)
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

func (c *SubstrateChain) Write(message *relayer.Message) error {
	return c.writer.VoteProposal(message)

}

func (c *SubstrateChain) ChainID() uint8 {
	return c.chainID
}

func (c *SubstrateChain) Stop() {
	close(c.stop)
}
