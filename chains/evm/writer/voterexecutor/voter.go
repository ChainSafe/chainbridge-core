package voterexecutor

import (
	"errors"
	"time"

	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
)

type VoterExecutor struct {
	stop     <-chan struct{}
	sysErr   chan<- error
	txSender TxSender
}

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Time between retrying a failed tx
const TxRetryLimit = 10

type TxSender interface {
	LockAndUpdateNonce()
	UnlockOpts()
	Opts() *bind.TransactOpts
}

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")

func (v *VoterExecutor) ExecuteProposal(proposal relayer.Proposal) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-v.stop:
			return
		default:
			// Nonce
			v.txSender.LockAndUpdateNonce()

			tx, err := v.bridgeContract.ExecuteProposal(
				v.txSender.Opts(),
				uint8(proposal.GetSource()),
				uint64(proposal.GetDepositNonce()),
				proposal.GetProposalData(),
				proposal.GetResourceID(),
			)
			v.txSender.UnlockOpts()

			if err == nil {
				log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
				return
			}
			if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				log.Error().Err(err).Msg("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
				time.Sleep(TxRetryInterval)
			}
			// Checking proposal status one more time (Since it could be execute by some other bridge). If it is finalized then we do not need to retry
			if proposal.ProposalIsComplete() {
				log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Proposal finalized on chain")
				return
			}
		}
	}
}
