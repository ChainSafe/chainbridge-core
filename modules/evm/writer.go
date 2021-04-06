package evm

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/Bridge"
	"github.com/ChainSafe/chainbridgev2/relayer"
	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type ContractCaller interface {
	FilterLogs(ctx context.Context, q goeth.FilterQuery) ([]types.Log, error)
	LatestBlock() (*big.Int, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	CallOpts() *bind.CallOpts
	Opts() *bind.TransactOpts
	LockAndUpdateOpts() error
	UnlockOpts()
	WaitForBlock(block *big.Int) error
}

// Number of blocks to wait for an finalization event
const ExecuteBlockWatchLimit = 100

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Maximum number of tx retries before exiting
const TxRetryLimit = 10

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")
var ErrFatalTx = errors.New("submission of transaction failed")
var ErrFatalQuery = errors.New("query of chain state failed")

type Bridger interface {
	GetProposal(opts *bind.CallOpts, originChainID uint8, depositNonce uint64, dataHash [32]byte) (Bridge.BridgeProposal, error)
	HasVotedOnProposal(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte, arg2 common.Address) (bool, error)
	VoteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, resourceID [32]byte, dataHash [32]byte) (*types.Transaction, error)
	ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte) (*types.Transaction, error)
}

type EVMWriter struct {
	chainID        uint8
	client         ContractCaller
	stop           <-chan struct{}
	sysErr         chan<- error
	bridgeContract Bridger
}

// NewWriter creates and returns writer
func NewWriter(client ContractCaller, bridgeContract Bridger, stop <-chan struct{}, sysErr chan<- error) *EVMWriter {
	return &EVMWriter{
		client: client,
		stop:   stop,
		sysErr: sysErr,
	}
}

func (w *EVMWriter) Write(m relayer.XCMessager) {
	data, err := m.CreateProposalData()
	if err != nil {
		panic(err)
	}
	dataHash := m.CreateProposalDataHash(data)
	//
	//if !w.shouldVote(m, dataHash) {
	//	if w.proposalIsPassed(m.Source, m.DepositNonce, dataHash) {
	//		// We should not vote for this proposal but it is ready to be executed
	//		w.executeProposal(m, data, dataHash)
	//		return true
	//	} else {
	//		return false
	//	}
	//}
	// Capture latest block so when know where to watch from
	//latestBlock, err := w.client.LatestBlock()
	//if err != nil {
	//	log.Error().Err(err).Msg("unable to fetch latest block")
	//	return false
	//}
	//
	//// watch for execution event
	//go w.watchThenExecute(m, data, dataHash, latestBlock)
	//
	w.voteProposal(m, dataHash)
}

// voteProposal submits a vote proposal
// a vote proposal will try to be submitted up to the TxRetryLimit times
func (w *EVMWriter) voteProposal(m relayer.XCMessager, dataHash ethcommon.Hash) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			// Checking first does proposal complete? If so, we do not need to vote for it
			//if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
			//	log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Proposal voting complete on chain")
			//	return
			//}
			err := w.client.LockAndUpdateOpts()
			if err != nil {
				log.Error().Err(err).Msg("Failed to update tx opts")
				continue
			}

			tx, err := w.bridgeContract.VoteProposal(
				w.client.Opts(),
				uint8(m.GetSource()),
				uint64(m.GetDepositNonce()),
				m.GetResourceID(),
				dataHash,
			)
			w.client.UnlockOpts()
			if err != nil {
				if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
					log.Debug().Msg("Nonce too low, will retry")
					time.Sleep(TxRetryInterval)
					continue
				} else {
					log.Warn().Interface("source", m.GetSource()).Interface("dest", m.GetDestination()).Interface("nonce", m.GetDepositNonce()).Msg("Voting failed")
					time.Sleep(TxRetryInterval)
					continue
				}
			}
			log.Info().Str("tx", tx.Hash().Hex()).Interface("src", m.GetSource()).Interface("depositNonce", m.GetDepositNonce()).Msg("Submitted proposal vote")
			return
		}
	}
	log.Error().Interface("source", m.GetSource()).Interface("dest", m.GetDestination()).Interface("nonce", m.GetDepositNonce()).Msg("Submission of Vote transaction failed")
	w.sysErr <- ErrFatalTx
}

//
//// executeProposal executes the proposal
//func (w *EVMWriter) executeProposal(m relayer.XCMessager, data []byte, dataHash ethcommon.Hash) {
//	for i := 0; i < TxRetryLimit; i++ {
//		select {
//		case <-w.stop:
//			return
//		default:
//			err := w.client.LockAndUpdateOpts()
//			if err != nil {
//				log.Error().Err(err).Msg("Failed to update nonce")
//				return
//			}
//
//			tx, err := w.bridgeContract.ExecuteProposal(
//				w.client.Opts(),
//				uint8(m.GetSource()),
//				uint64(m.GetDepositNonce()),
//				data,
//				m.GetResourceID(),
//			)
//			w.client.UnlockOpts()
//
//			if err == nil {
//				log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
//				return
//			}
//			if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
//				log.Error().Err(err).Msg("Nonce too low, will retry")
//				time.Sleep(TxRetryInterval)
//			} else {
//				log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
//				time.Sleep(TxRetryInterval)
//			}
//			// Checking proposal status one more time (Since it could be execute by some other bridge). If it is finalized then we do not need to retry
//			if w.proposalIsFinalized(m.Source, m.DepositNonce, dataHash) {
//				log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Proposal finalized on chain")
//				return
//			}
//		}
//	}
//	log.Error().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Submission of Execute transaction failed")
//	w.sysErr <- ErrFatalTx
//}
//
//// watchThenExecute watches for the latest block and executes once the matching finalized event is found
//func (w *EVMWriter) watchThenExecute(m *utils.Message, data []byte, dataHash ethcommon.Hash, latestBlock *big.Int) {
//	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Watching for finalization event")
//
//	// watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
//	for i := 0; i < ExecuteBlockWatchLimit; i++ {
//		select {
//		case <-w.stop:
//			return
//		default:
//			// watch for the lastest block, retry up to BlockRetryLimit times
//			for waitRetrys := 0; waitRetrys <= BlockRetryLimit; waitRetrys++ {
//				err := w.client.WaitForBlock(latestBlock)
//				if err != nil {
//					log.Error().Err(err).Msg("Waiting for block failed")
//					// Exit if retries exceeded
//					if waitRetrys == BlockRetryLimit {
//						log.Error().Err(err).Msg("Waiting for block retries exceeded, shutting down")
//						w.sysErr <- ErrFatalQuery
//						return
//					}
//				} else {
//					break
//				}
//			}
//
//			// query for logs
//			query := buildQuery(w.cfg.BridgeContract, utils.ProposalEvent, latestBlock, latestBlock)
//			evts, err := w.client.FilterLogs(context.Background(), query)
//			if err != nil {
//				log.Error().Err(err).Msg("Failed to fetch logs")
//				return
//			}
//
//			// execute the proposal once we find the matching finalized event
//			for _, evt := range evts {
//				sourceId := evt.Topics[1].Big().Uint64()
//				depositNonce := evt.Topics[2].Big().Uint64()
//				status := evt.Topics[3].Big().Uint64()
//
//				if m.Source == utils.ChainId(sourceId) &&
//					m.DepositNonce.Big().Uint64() == depositNonce &&
//					utils.IsPassed(uint8(status)) {
//					w.executeProposal(m, data, dataHash)
//					return
//				} else {
//					log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Uint64("status", status).Msg("Ignoring event")
//				}
//			}
//			log.Trace().Interface("block", latestBlock).Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("No finalization event found in current block")
//			latestBlock = latestBlock.Add(latestBlock, big.NewInt(1))
//		}
//	}
//	log.Warn().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Block watch limit exceeded, skipping execution")
//}
