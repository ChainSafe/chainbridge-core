package relayer

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

// Number of blocks to wait for an finalization event
const ExecuteBlockWatchLimit = 100

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Time between retrying a failed tx
const TxRetryLimit = 10

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")
var ErrFatalTx = errors.New("submission of transaction failed")
var ErrFatalQuery = errors.New("query of chain state failed")

type ProposalVoter interface {
	VoteProposal(opts *bind.TransactOpts, m XCMessager, dataHash [32]byte) (*types.Transaction, error)
}

type ProposalExecuter interface {
	ExecuteProposal(m XCMessager, data []byte, dataHash ethcommon.Hash)
}

type ContractCaller interface {
	FilterLogs(ctx context.Context, q goeth.FilterQuery) ([]types.Log, error)
	LatestBlock() (*big.Int, error)
	ProposalStatusDefiner() ProposalStatus
	//BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	//CallOpts() *bind.CallOpts
	//Opts() *bind.TransactOpts
	//LockAndUpdateOpts() error
	//UnlockOpts()
	//WaitForBlock(block *big.Int) error
}

type Writer struct {
	voter        ProposalVoter
	executer     ProposalExecuter
	transactOpts *bind.TransactOpts
	stop         <-chan struct{}
	sysErr       chan<- error
	client       ContractCaller
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
func (w *Writer) watchThenExecute(m XCMessager, data []byte, dataHash ethcommon.Hash) {
	delay := big.NewInt(10)
	blockToWatch, err := w.client.LatestBlock()
	if err != nil {
		log.Error().Err(err).Msg("unable to fetch latest block")
		w.sysErr <- errors.New("unable to fetch latest block")
		return
	}
	log.Info().Interface("src", m.GetSource()).Interface("nonce", m.GetDepositNonce()).Msg("Watching for finalization event")

	// watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			// Waits chain to overtake current block for delay
			err := BlockWaiter(w.client, blockToWatch, delay, w.stop)
			if err != nil {
				w.sysErr <- err
				return
			}

			// query for logs
			query := buildQuery(w.cfg.BridgeContract, crypto.Keccak256Hash([]byte("ProposalEvent(uint8,uint64,uint8,bytes32,bytes32)")), blockToWatch, blockToWatch.Add(blockToWatch, delay))
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

				if m.GetSource() == sourceId &&
					m.GetDepositNonce() == depositNonce &&
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

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig ethcommon.Hash, startBlock *big.Int, endBlock *big.Int) goeth.FilterQuery {
	query := goeth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []ethcommon.Address{contract},
		Topics: [][]ethcommon.Hash{
			{sig},
		},
	}
	return query
}
