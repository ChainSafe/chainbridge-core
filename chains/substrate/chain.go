package substrate

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/store"
	"github.com/rs/zerolog/log"
)

type ProposalVoter interface {
	VoteProposal(message *message.Message) error
}

type EventListener interface {
	ListenToEvents(startBlock *big.Int, domainID uint8, blockstore *store.BlockStore, stopChn <-chan struct{}, errChn chan<- error) <-chan *message.Message
}

type SubstrateChain struct {
	domainID   uint8
	stop       chan<- struct{}
	listener   EventListener
	writer     ProposalVoter
	blockstore *store.BlockStore
	config     *chain.SubstrateConfig
}

func NewSubstrateChain(listener EventListener, writer ProposalVoter, blockstore *store.BlockStore, domainID uint8, config *chain.SubstrateConfig) *SubstrateChain {
	return &SubstrateChain{
		listener:   listener,
		writer:     writer,
		blockstore: blockstore,
		domainID:   domainID,
		config:     config,
	}
}

func (c *SubstrateChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan *message.Message) {
	log.Info().Msg("Polling Blocks...")

	startingBlock, err := c.blockstore.GetStartBlock(
		*c.config.GeneralChainConfig.Id,
		c.config.StartBlock,
		c.config.GeneralChainConfig.LatestBlock,
		c.config.GeneralChainConfig.FreshStart,
	)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}

	ech := c.listener.ListenToEvents(startingBlock, c.domainID, c.blockstore, stop, sysErr)
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
