package client

import (
	"context"
	"math/big"
	"sync"

	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

func NewClient(endpoint string, http bool, kp *secp256k1.Keypair, gasLimit *big.Int, gasPrice *big.Int) (*Client, error) {
	c := &Client{
		endpoint:    endpoint,
		http:        http,
		kp:          kp,
		maxGasPrice: gasPrice,
		gasLimit:    gasLimit,
		stop:        make(chan int),
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}
	return c, nil
}

type Client struct {
	*ethclient.Client
	endpoint    string
	http        bool
	kp          *secp256k1.Keypair
	gasLimit    *big.Int
	maxGasPrice *big.Int
	opts        *bind.TransactOpts
	callOpts    *bind.CallOpts
	nonce       uint64
	nonceLock   sync.Mutex
	optsLock    sync.Mutex
	stop        chan int // All routines should exit when this channel is closed
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
func (c *Client) Connect() error {
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
	c.nonce = 0
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

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, id) // TODO pass ChainID thru config somehow
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
