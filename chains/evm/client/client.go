package client

import (
	"context"
	"fmt"
	"math/big"

	bridgeHandler "github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/Bridge"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

func NewClient(endpoint string, http bool, stopChan <-chan struct{}, errChan chan<- error, bridgeAddr common.Address) (*Client, error) {
	c := &Client{
		endpoint: endpoint,
		http:     http,
		stopChn:  stopChan,
		errChn:   errChan,
	}
	if err := c.connect(); err != nil {
		return nil, err
	}

	bridgeContract, err := bridgeHandler.NewBridge(bridgeAddr, c)
	if err != nil {
		return nil, err
	}
	c.bridgeContract = bridgeContract
	return c, nil
}

type Client struct {
	client         *ethclient.Client
	endpoint       string
	http           bool
	bridgeContract *bridgeHandler.Bridge
	stopChn        <-chan struct{}
	errChn         chan<- error
}

// LatestBlock returns the latest block from the current chain
func (c *Client) LatestBlock() (*big.Int, error) {
	header, err := c.client.HeaderByNumber(context.Background(), nil)
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
	c.client.Client = ethclient.NewClient(rpcClient)

	// Construct tx opts, call opts, and nonce mechanism
	//opts, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice)
	//if err != nil {
	//	return err
	//}
	//c.opts = opts
	//c.callOpts = &bind.CallOpts{From: c.senderKP.CommonAddress()}
	return nil
}

// This is done because we can't pass ethclient.Client in to the evm.Writer because client.ChainID function is not a part of any interfaces inside go-ethereum library
func (c *Client) GetEthClient() *ethclient.Client {
	return c.client
}

func (c *Client) MatchResourceIDToHandlerAddress(rID [32]byte) (string, error) {
	addr, err := c.bridgeContract.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
	if err != nil {
		return "", fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rID, err)
	}
	return addr.String(), nil
}
