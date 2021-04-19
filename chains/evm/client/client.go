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

func NewClient(endpoint string, http bool) (*Client, error) {
	c := &Client{
		endpoint: endpoint,
		http:     http,
	}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

type Client struct {
	*ethclient.Client
	endpoint string
	http     bool
	stopChn  <-chan struct{}
	errChn   chan<- error
}

// LatestBlock returns the latest block from the current chain
func (c *Client) LatestBlock() (*big.Int, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), nil)
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
	return nil
}

func (c *Client) GetEthClient() *ethclient.Client {
	return c.Client
}

func (c *Client) MatchResourceIDToHandlerAddress(bridgeAddr common.Address, rID [32]byte) (string, error) {
	bridgeContract, err := bridgeHandler.NewBridge(bridgeAddr, c)
	if err != nil {
		return "", err
	}
	addr, err := bridgeContract.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
	if err != nil {
		return "", fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rID, err)
	}
	return addr.String(), nil
}
