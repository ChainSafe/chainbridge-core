package relayer

import (
	"context"
	"errors"
	"math/big"
	"time"

	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

// Number of blocks to wait for an finalization event
const ExecuteBlockWatchLimit = 100

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Time between retrying a failed tx
const TxRetryLimit = 10

//var ErrNonceTooLow = errors.New("nonce too low")
//var ErrTxUnderpriced = errors.New("replacement transaction underpriced")
//var ErrFatalTx = errors.New("submission of transaction failed")
//var ErrFatalQuery = errors.New("query of chain state failed")

type ProposalVoter interface {
	VoteProposal(proposal Proposal) (*types.Transaction, error)
}

type ProposalExecuter interface {
	ExecuteProposal(proposal Proposal)
}

type ContractCaller interface {
	FilterLogs(ctx context.Context, q goeth.FilterQuery) ([]types.Log, error)
	LatestBlock() (*big.Int, error)
}

type Writer struct {
	voter          ProposalVoter
	executer       ProposalExecuter
	transactOpts   *bind.TransactOpts
	stop           <-chan struct{}
	sysErr         chan<- error
	client         ContractCaller
	propCreatorFn  ProposalCreatorFn
	BridgeContract ethcommon.Address
}

func NewWriter(voter ProposalVoter, executer ProposalExecuter, transactOpts *bind.TransactOpts, propCreatorFn ProposalCreatorFn) *Writer {
	return &Writer{
		voter:         voter,
		executer:      executer,
		transactOpts:  transactOpts,
		propCreatorFn: propCreatorFn,
	}
}

func (w *Writer) Write(m XCMessager) {

	prop, err := w.propCreatorFn(m)
	if err != nil {
		w.sysErr <- err
	}

	if !prop.ShouldVoteFor() {
		if prop.GetProposalStatus() == ProposalStatusPassed {
			// We should not vote for this proposal but it is ready to be executed
			w.executer.ExecuteProposal(prop)
			return
		} else {
			return
		}
	}

	go w.watchThenExecute(prop)
	w.voter.VoteProposal(prop)
}

// watchThenExecute watches for the latest block and executes once the matching finalized event is found
func (w *Writer) watchThenExecute(prop Proposal) {
	delay := big.NewInt(10)
	blockToWatch, err := w.client.LatestBlock()
	if err != nil {
		log.Error().Err(err).Msg("unable to fetch latest block")
		w.sysErr <- errors.New("unable to fetch latest block")
		return
	}
	log.Info().Interface("src", prop.GetSource()).Interface("nonce", prop.GetDepositNonce()).Msg("Watching for finalization event")

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
			query := buildQuery(w.BridgeContract, crypto.Keccak256Hash([]byte("ProposalEvent(uint8,uint64,uint8,bytes32,bytes32)")), blockToWatch, blockToWatch.Add(blockToWatch, delay))
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

				if prop.GetSource() == uint8(sourceId) &&
					prop.GetDepositNonce() == depositNonce &&
					prop.GetProposalStatus() == ProposalStatusPassed {
					w.executer.ExecuteProposal(prop)
					return
				} else {
					log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Uint64("status", status).Msg("Ignoring event")
				}
			}
			log.Trace().Interface("block", blockToWatch).Interface("src", prop.GetSource()).Interface("nonce", prop.GetDepositNonce()).Msg("No finalization event found in current block")
			blockToWatch = blockToWatch.Add(blockToWatch, big.NewInt(1))
		}
	}
	log.Warn().Interface("source", prop.GetSource()).Interface("dest", prop.GetDestination()).Interface("nonce", prop.GetDepositNonce()).Msg("Block watch limit exceeded, skipping execution")
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
