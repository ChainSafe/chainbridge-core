// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package client

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
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

var ErrFatalTx = errors.New("submission of transaction failed")
var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Tries to retry sending transaction
const TxRetryLimit = 10

func NewEVMClient(endpoint string, http bool, sender *secp256k1.Keypair) (*EVMClient, error) {
	c := &EVMClient{
		endpoint: endpoint,
		http:     http,
		sender:   sender,
	}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

type EVMClient struct {
	*ethclient.Client
	endpoint      string
	http          bool
	stop          <-chan struct{}
	errChn        chan<- error
	optsLock      sync.Mutex
	opts          *bind.TransactOpts
	sender        *secp256k1.Keypair
	maxGasPrice   *big.Int   // TODO
	gasMultiplier *big.Float // TODO
}

// Connect starts the ethereum WS connection
func (c *EVMClient) connect() error {
	log.Info().Str("url", c.endpoint).Msg("Connecting to ethereum chain...")
	var rpcClient *rpc.Client
	var err error
	// Start http or ws client
	if c.http {
		rpcClient, err = rpc.DialHTTP(c.endpoint)
	} else {
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.endpoint, "/ws")
	}
	if err != nil {
		return err
	}
	c.Client = ethclient.NewClient(rpcClient)
	return nil
}

// LatestBlock returns the latest block from the current chain
func (c *EVMClient) LatestBlock() (*big.Int, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

func (c *EVMClient) GetEthClient() *ethclient.Client {
	return c.Client
}

func (c *EVMClient) ExecuteProposal(bridgeAddress string, proposal relayer.Proposal) error {
	for i := 0; i < TxRetryLimit; i++ {
		err := c.lockAndUpdateOpts()
		if err != nil {
			log.Error().Err(err).Msgf("failed to update tx opts")
			time.Sleep(TxRetryInterval)
		}
		b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
		if err != nil {
			return err
		}
		tx, err := b.ExecuteProposal(
			c.getOpts(),
			uint8(proposal.GetSource()),
			uint64(proposal.GetDepositNonce()),
			proposal.GetProposalData(),
			proposal.GetResourceID(),
		)
		c.unlockOpts()
		if err == nil {
			log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
			return nil
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
		s, err := c.ProposalStatus(bridgeAddress, proposal)
		if err != nil {
			log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
			continue
		}
		if s == relayer.ProposalStatusPassed || s == relayer.ProposalStatusExecuted || s == relayer.ProposalStatusCanceled {
			log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Proposal finalized on chain")
			return nil
		}
	}
	log.Error().Msgf("Submission of Execution transaction failed, source %v dest %v depNonce %v", proposal.GetSource(), proposal.GetDestination(), proposal.GetDepositNonce())
	return ErrFatalTx
}

func (c *EVMClient) VoteProposal(bridgeAddress string, proposal relayer.Proposal) error {
	for i := 0; i < TxRetryLimit; i++ {
		err := c.lockAndUpdateOpts()
		if err != nil {
			log.Error().Err(err).Msgf("failed to update tx opts")
		}
		b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
		if err != nil {
			return err
		}
		tx, err := b.VoteProposal(
			c.getOpts(),
			uint8(proposal.GetSource()),
			uint64(proposal.GetDepositNonce()),
			proposal.GetResourceID(),
			proposal.GetProposalDataHash(),
		)
		c.unlockOpts()
		if err == nil {
			log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal vote")
			return nil
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
		ps, err := c.ProposalStatus(bridgeAddress, proposal)
		if err != nil {
			log.Error().Err(err).Msgf("error getting proposal status %+v", proposal)
			continue
		}
		if ps == relayer.ProposalStatusPassed {
			log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Proposal is ready to be executed on chain")
			return nil
		}
	}
	log.Error().Msgf("Submission of vote transaction failed, source %v dest %v depNonce %v", proposal.GetSource(), proposal.GetDestination(), proposal.GetDepositNonce())
	return ErrFatalTx
}

func (c *EVMClient) ProposalStatus(bridgeAddress string, p relayer.Proposal) (relayer.ProposalStatus, error) {
	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
	if err != nil {
		return 99, err
	}
	prop, err := b.GetProposal(&bind.CallOpts{}, p.GetSource(), p.GetDepositNonce(), p.GetProposalDataHash())
	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return 99, err
	}
	return relayer.ProposalStatus(prop.Status), nil
}

func (c *EVMClient) VotedBy(bridgeAddress string, p relayer.Proposal) bool {
	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
	if err != nil {
		return false
	}
	hv, err := b.HasVotedOnProposal(&bind.CallOpts{}, p.GetIDAndNonce(), p.GetProposalDataHash(), c.sender.CommonAddress())
	if err != nil {
		return false
	}
	return hv
}

func (c *EVMClient) MatchResourceIDToHandlerAddress(bridgeAddress string, rID [32]byte) (string, error) {
	b, err := Bridge.NewBridge(common.HexToAddress(bridgeAddress), c)
	if err != nil {
		return "", err
	}
	addr, err := b.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
	if err != nil {
		return "", fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rID, err)
	}
	return addr.String(), nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *EVMClient) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, error) {
	privateKey := c.sender.PrivateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	id, err := c.ChainID(context.Background())
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

func (c *EVMClient) unlockOpts() {
	c.optsLock.Unlock()
}

func (c *EVMClient) lockAndUpdateOpts() error {
	c.optsLock.Lock()

	gasPrice, err := c.safeEstimateGas(context.TODO())
	if err != nil {
		return err
	}
	c.opts.GasPrice = gasPrice

	nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *EVMClient) getOpts() *bind.TransactOpts {
	return c.opts
}

func (c *EVMClient) safeEstimateGas(ctx context.Context) (*big.Int, error) {
	suggestedGasPrice, err := c.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	gasPrice := multiplyGasPrice(suggestedGasPrice, c.gasMultiplier)

	// Check we aren't exceeding our limit
	if gasPrice.Cmp(c.maxGasPrice) == 1 {
		return c.maxGasPrice, nil
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
