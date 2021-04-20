package bridger

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/Bridge"
	"github.com/ChainSafe/chainbridgev2/crypto/secp256k1"
	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type EVMBridgeClient struct {
	stop           <-chan struct{}
	errChn         chan<- error
	bridgeContract *Bridge.Bridge
	optsLock       sync.Mutex
	opts           *bind.TransactOpts
	sender         *secp256k1.Keypair
	maxGasPrice    *big.Int   // TODO
	gasMultiplier  *big.Float // TODO
	backend        *ethclient.Client
}

var ErrFatalTx = errors.New("submission of transaction failed")

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Tries to retry sending transaction
const TxRetryLimit = 10

func NewBridgeClient(bridgeAddress common.Address, backend *ethclient.Client, sender *secp256k1.Keypair, stop <-chan struct{}, errChn chan<- error) (*EVMBridgeClient, error) {
	b, err := Bridge.NewBridge(bridgeAddress, backend)
	if err != nil {
		return nil, err
	}
	return &EVMBridgeClient{
		stop:           stop,
		errChn:         errChn,
		bridgeContract: b,
		backend:        backend,
		sender:         sender,
	}, nil
}

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")

func (v *EVMBridgeClient) ExecuteProposal(proposal relayer.Proposal) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-v.stop:
			return
		default:
			err := v.lockAndUpdateOpts()
			if err != nil {
				log.Error().Err(err).Msgf("failed to update tx opts")
				time.Sleep(TxRetryInterval)
			}
			tx, err := v.bridgeContract.ExecuteProposal(
				v.getOpts(),
				uint8(proposal.GetSource()),
				uint64(proposal.GetDepositNonce()),
				proposal.GetProposalData(),
				proposal.GetResourceID(),
			)
			v.unlockOpts()
			if err == nil {
				log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
				return
			}
			if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				log.Error().Err(err).Msg("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				// TODO: this part is unclear. Does sending transaction with contract binding response with error if transaction failed inside contract?
				log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
				time.Sleep(TxRetryInterval)
			}
			// Checking proposal status one more time (Since it could be execute by some other bridge). If it is completed then we do not need to retry
			s, err := v.ProposalStatus(proposal)
			if err != nil {
				log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
				continue
			}
			if s == relayer.ProposalStatusPassed || s == relayer.ProposalStatusExecuted || s == relayer.ProposalStatusCanceled {
				log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Proposal finalized on chain")
				return
			}
		}
	}
	log.Error().Msgf("Submission of Execution transaction failed, source %v dest %v depNonce %v", proposal.GetSource(), proposal.GetDestination(), proposal.GetDepositNonce())
	v.errChn <- ErrFatalTx
}

func (v *EVMBridgeClient) VoteProposal(proposal relayer.Proposal) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-v.stop:
			return
		default:
			err := v.lockAndUpdateOpts()
			if err != nil {
				log.Error().Err(err).Msgf("failed to update tx opts")
			}
			tx, err := v.bridgeContract.VoteProposal(
				v.getOpts(),
				uint8(proposal.GetSource()),
				uint64(proposal.GetDepositNonce()),
				proposal.GetResourceID(),
				proposal.GetProposalDataHash(),
			)
			v.unlockOpts()
			if err == nil {
				log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal vote")
				return
			}
			if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				log.Error().Err(err).Msg("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				// TODO: this part is unclear. Does sending transaction with contract binding response with error if transaction failed inside contract?
				log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
				time.Sleep(TxRetryInterval)
			}
			// Checking proposal status one more time (Since it could be execute by some other bridge). If it is completed then we do not need to retry
			ps, err := v.ProposalStatus(proposal)
			if err != nil {
				log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
				continue
			}
			if ps == relayer.ProposalStatusPassed {
				log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Proposal is ready to be executed on chain")
				return
			}
		}
	}
	log.Error().Msgf("Submission of vote transaction failed, source %v dest %v depNonce %v", proposal.GetSource(), proposal.GetDestination(), proposal.GetDepositNonce())
	v.errChn <- ErrFatalTx
}

func (v *EVMBridgeClient) ProposalStatus(p relayer.Proposal) (relayer.ProposalStatus, error) {
	prop, err := v.bridgeContract.GetProposal(&bind.CallOpts{}, p.GetSource(), p.GetDepositNonce(), p.GetProposalDataHash())
	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return 999, err
	}
	return relayer.ProposalStatus(prop.Status), nil
}

func (v *EVMBridgeClient) VotedBy(p relayer.Proposal) bool {
	b, err := v.bridgeContract.HasVotedOnProposal(&bind.CallOpts{}, p.GetIDAndNonce(), p.GetProposalDataHash(), v.sender.CommonAddress())
	if err != nil {
		return false
	}
	return b
}

func (v *EVMBridgeClient) MatchResourceIDToHandlerAddress(rID [32]byte) (string, error) {
	addr, err := v.bridgeContract.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
	if err != nil {
		return "", fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rID, err)
	}
	return addr.String(), nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (v *EVMBridgeClient) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, error) {
	privateKey := v.sender.PrivateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := v.backend.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	id, err := v.backend.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, id)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	return auth, nil
}

func (v *EVMBridgeClient) unlockOpts() {
	v.optsLock.Unlock()
}

func (v *EVMBridgeClient) lockAndUpdateOpts() error {
	v.optsLock.Lock()

	gasPrice, err := v.safeEstimateGas(context.TODO())
	if err != nil {
		return err
	}
	v.opts.GasPrice = gasPrice

	nonce, err := v.backend.PendingNonceAt(context.Background(), v.opts.From)
	if err != nil {
		v.optsLock.Unlock()
		return err
	}
	v.opts.Nonce.SetUint64(nonce)
	return nil
}

func (v *EVMBridgeClient) getOpts() *bind.TransactOpts {
	return v.opts
}

func (v *EVMBridgeClient) safeEstimateGas(ctx context.Context) (*big.Int, error) {
	suggestedGasPrice, err := v.backend.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	gasPrice := multiplyGasPrice(suggestedGasPrice, v.gasMultiplier)

	// Check we aren't exceeding our limit
	if gasPrice.Cmp(v.maxGasPrice) == 1 {
		return v.maxGasPrice, nil
	} else {
		return gasPrice, nil
	}
}

func multiplyGasPrice(gasEstimate *big.Int, gasMultiplier *big.Float) *big.Int {
	gasEstimateFloat := new(big.Float).SetInt(gasEstimate)
	result := gasEstimateFloat.Mul(gasEstimateFloat, gasMultiplier)
	gasPrice := new(big.Int)
	result.Int(gasPrice)
	return gasPrice
}
