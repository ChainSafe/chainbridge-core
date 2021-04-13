package evmd

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/Bridge"
	"github.com/ChainSafe/chainbridgev2/relayer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

type Bridger interface {
	GetProposal(opts *bind.CallOpts, originChainID uint8, depositNonce uint64, dataHash [32]byte) (Bridge.BridgeProposal, error)
	HasVotedOnProposal(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte, arg2 common.Address) (bool, error)
	VoteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, resourceID [32]byte, dataHash [32]byte) (*types.Transaction, error)
	ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte, signatureHeader []byte, aggregatePublicKey []byte, hashedMessage [32]byte, rootHash [32]byte, key []byte, nodes []byte) (*types.Transaction, error)
}

func NewClient(endpoint string, http bool, kp *secp256k1.Keypair, gasLimit *big.Int, gasPrice *big.Int, stopChan <-chan struct{}, errChan chan<- error) (*Client, error) {
	c := &Client{
		endpoint:    endpoint,
		http:        http,
		kp:          kp,
		maxGasPrice: gasPrice,
		gasLimit:    gasLimit,
		stopChn:     stopChan,
		errChn:      errChan,
	}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

type Client struct {
	*ethclient.Client
	endpoint       string
	http           bool
	kp             *secp256k1.Keypair
	gasLimit       *big.Int
	maxGasPrice    *big.Int
	opts           *bind.TransactOpts
	callOpts       *bind.CallOpts
	nonceLock      sync.Mutex
	optsLock       sync.Mutex
	stopChn        <-chan struct{}
	bridgeContract Bridger
	errChn         chan<- error
}

// LatestBlock returns the latest block from the current chain
func (c *Client) LatestBlock() (*big.Int, error) {
	header, err := c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// Connect starts the ethereum WS connection
func (c *Client) connect() error {
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

	// Construct tx opts, call opts, and nonce mechanism
	opts, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice)
	if err != nil {
		return err
	}
	c.opts = opts
	c.callOpts = &bind.CallOpts{From: c.kp.CommonAddress()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Client) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

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

func (c *Client) unlockNonce() {
	c.nonceLock.Unlock()
}

func (c *Client) unlockOpts() {
	c.optsLock.Unlock()
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Client) lockAndUpdateOpts() error {
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

func (c *Client) lockAndUpdateNonce() error {
	c.nonceLock.Lock()
	nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.nonceLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

// This should be done as function that accepts Config
//func (c *Client) safeEstimateGas(ctx context.Context) (*big.Int, error) {
//	suggestedGasPrice, err := c.SuggestGasPrice(context.TODO())
//	if err != nil {
//		return nil, err
//	}
//
//	gasPrice := multiplyGasPrice(suggestedGasPrice, c.gasMultiplier)
//	// Check we aren't exceeding our limit
//	if gasPrice.Cmp(c.maxGasPrice) == 1 {
//		return c.maxGasPrice, nil
//	} else {
//		return gasPrice, nil
//	}
//}

func multiplyGasPrice(gasEstimate *big.Int, gasMultiplier *big.Float) *big.Int {
	gasEstimateFloat := new(big.Float).SetInt(gasEstimate)

	result := gasEstimateFloat.Mul(gasEstimateFloat, gasMultiplier)

	gasPrice := new(big.Int)

	result.Int(gasPrice)

	return gasPrice
}

// Maximum number of tx retries before exiting
const TxRetryLimit = 10
const TxRetryInterval = time.Second * 2

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")
var ErrFatalTx = errors.New("submission of transaction failed")
var ErrFatalQuery = errors.New("query of chain state failed")

func (c *Client) VoteProposal(proposal relayer.Proposal) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-c.stopChn:
			return
		default:
			// Checking first does proposal complete? If so, we do not need to vote for it
			if relayer.ProposalIsComplete(proposal) {
				log.Info().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Proposal voting complete on chain")
				return
			}
			err := c.lockAndUpdateOpts()
			if err != nil {
				log.Error().Err(err).Msg("Failed to update tx opts")
				continue
			}

			tx, err := c.bridgeContract.VoteProposal(
				c.opts,
				proposal.GetSource(),
				proposal.GetDepositNonce(),
				proposal.GetResourceID(),
				proposal.GetProposalDataHash(proposal.GetProposalData()),
			)
			c.unlockOpts()
			if err != nil {
				if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
					log.Debug().Msg("Nonce too low, will retry")
					time.Sleep(TxRetryInterval)
					continue
				} else {
					log.Warn().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Voting failed")
					time.Sleep(TxRetryInterval)
					continue
				}
			}
			log.Info().Str("tx", tx.Hash().Hex()).Interface("src", proposal.GetSource()).Interface("depositNonce", proposal.GetDepositNonce()).Msg("Submitted proposal vote")
			return
		}
	}
	log.Error().Interface("source", proposal.GetSource()).Interface("dest", proposal.GetDestination()).Interface("nonce", proposal.GetDepositNonce()).Msg("Submission of Vote transaction failed")
	c.errChn <- ErrFatalTx
}
