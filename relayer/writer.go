package relayer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type ProposalVoter interface {
	VoteProposal(opts *bind.TransactOpts, m XCMessager, dataHash [32]byte) (*types.Transaction, error)
}

type ProposalExecuter interface {
	ExecuteProposal(m XCMessager, data []byte, dataHash ethcommon.Hash)
}

type Writer struct {
	voter        ProposalVoter
	executer     ProposalExecuter
	transactOpts *bind.TransactOpts
}

func NewWriter(voter ProposalVoter, executer ProposalExecuter, transactOpts *bind.TransactOpts) *Writer {
	return &Writer{
		voter:        voter,
		executer:     executer,
		transactOpts: transactOpts,
	}
}

func (w *Writer) Write(m XCMessager) {
	data, err := m.CreateProposalData()
	if err != nil {
		panic(err)
	}
	dataHash := m.CreateProposalDataHash(data)

	//if !w.shouldVote(m, dataHash) {
	//	if w.proposalIsPassed(m.Source, m.DepositNonce, dataHash) {
	//		// We should not vote for this proposal but it is ready to be executed
	//		w.executeProposal(m, data, dataHash)
	//		return true
	//	} else {
	//		return false
	//	}
	//}

	w.voter.VoteProposal(transactOpts, m, dataHash)

	w.watchThenExecute(query)
}

// watchThenExecute watches for the latest block and executes once the matching finalized event is found
func (w *Writer) watchThenExecute(m *utils.Message, data []byte, dataHash ethcommon.Hash, latestBlock *big.Int) {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Watching for finalization event")

	// watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			// watch for the lastest block, retry up to BlockRetryLimit times
			for waitRetrys := 0; waitRetrys <= BlockRetryLimit; waitRetrys++ {
				err := w.client.WaitForBlock(latestBlock)
				if err != nil {
					log.Error().Err(err).Msg("Waiting for block failed")
					// Exit if retries exceeded
					if waitRetrys == BlockRetryLimit {
						log.Error().Err(err).Msg("Waiting for block retries exceeded, shutting down")
						w.sysErr <- ErrFatalQuery
						return
					}
				} else {
					break
				}
			}

			// query for logs
			query := buildQuery(w.cfg.BridgeContract, utils.ProposalEvent, latestBlock, latestBlock)
			evts, err := w.client.FilterLogs(context.Background(), query)
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch logs")
				return
			}

			// execute the proposal once we find the matching finalized event
			for _, evt := range evts {
				sourceId := evt.Topics[1].Big().Uint64()
				depositNonce := evt.Topics[2].Big().Uint64()
				status := evt.Topics[3].Big().Uint64()

				if m.Source == utils.ChainId(sourceId) &&
					m.DepositNonce.Big().Uint64() == depositNonce &&
					utils.IsPassed(uint8(status)) {
					w.executer.ExecuteProposal(m, data, dataHash)
					return
				} else {
					log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Uint64("status", status).Msg("Ignoring event")
				}
			}
			log.Trace().Interface("block", latestBlock).Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("No finalization event found in current block")
			latestBlock = latestBlock.Add(latestBlock, big.NewInt(1))
		}
	}
	log.Warn().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Block watch limit exceeded, skipping execution")
}
