package substrate

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/rs/zerolog/log"
)

type ProposalVoter interface {
	VoteProposal(message *message.Message) error
}

type EventListener interface {
	ListenToEvents(startBlock *big.Int, domainID uint8, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan *message.Message
}

type SubstrateChain struct {
	domainID uint8
	stop     chan<- struct{}
	listener EventListener
	writer   ProposalVoter
	kvdb     blockstore.KeyValueReaderWriter
	config   *chain.SharedSubstrateConfig
}

func NewSubstrateChain(listener EventListener, writer ProposalVoter, kvdb blockstore.KeyValueReaderWriter, domainID uint8, config *chain.SharedSubstrateConfig) *SubstrateChain {
	return &SubstrateChain{
		listener: listener,
		writer:   writer,
		kvdb:     kvdb,
		domainID: domainID,
		config:   config,
	}
}

func (c *SubstrateChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan *message.Message) {
	log.Info().Msg("Polling Blocks...")
	// Handler chain specific configs and flags
	//b, err := blockstore.GetLastStoredBlock(c.kvdb, c.domainID)
	block, err := blockstore.SetupBlockstore(&c.config.GeneralChainConfig, c.kvdb, c.config.StartBlock)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}
	ech := c.listener.ListenToEvents(block, c.domainID, c.kvdb, stop, sysErr)
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

func (c *SubstrateChain) Write(message *message.Message) error {
	return c.writer.VoteProposal(message)

}

func (c *SubstrateChain) DomainID() uint8 {
	return c.domainID
}

func (c *SubstrateChain) Stop() {
	close(c.stop)
}
