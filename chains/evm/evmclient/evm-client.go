package evmclient

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

type config struct {
	maxGasPrice    *big.Int
	gasMultiplier  *big.Float
	relayerAddress common.Address
	kp             *secp256k1.Keypair
}

type EVMClient struct {
	*ethclient.Client
	rpClient  *rpc.Client
	nonceLock sync.Mutex
	config    *config
	nonce     *big.Int
}

type CommonTransaction interface {
	// Hash returns the transaction hash.
	Hash() common.Hash
	// Returns signed transaction by provided private key
	RawWithSignature(key *ecdsa.PrivateKey, chainID *big.Int) ([]byte, error)
}

func NewEVMClient(endpoint string, kp *secp256k1.Keypair) (*EVMClient, error) {
	log.Info().Str("url", endpoint).Msg("Connecting to evm chain...")
	rpcClient, err := rpc.DialContext(context.TODO(), endpoint)
	if err != nil {
		return nil, err
	}
	c := &config{
		kp: kp,
	}
	return &EVMClient{
		Client:   ethclient.NewClient(rpcClient),
		rpClient: rpcClient,
		config:   c,
	}, nil
}

// TO implement interface Configurable
func (c *EVMClient) Configurate() {
	c.config.maxGasPrice = big.NewInt(20000000000)
	c.config.gasMultiplier = big.NewFloat(1)

}

// LatestBlock returns the latest block from the current chain
func (c *EVMClient) LatestBlock() (*big.Int, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

const (
	DepositSignature string = "Deposit(uint8,bytes32,uint64)"
)

func (c *EVMClient) FetchDepositLogs(ctx context.Context, contractAddress common.Address, startBlock *big.Int, endBlock *big.Int) ([]*listener.DepositLogs, error) {
	logs, err := c.FilterLogs(ctx, buildQuery(contractAddress, DepositSignature, startBlock, endBlock))
	if err != nil {
		return nil, err
	}
	depositLogs := make([]*listener.DepositLogs, 0)

	for _, l := range logs {
		dl := &listener.DepositLogs{
			DestinationID: uint8(l.Topics[1].Big().Uint64()),
			ResourceID:    l.Topics[2],
			DepositNonce:  l.Topics[3].Big().Uint64(),
		}
		depositLogs = append(depositLogs, dl)
	}
	return depositLogs, nil
}

// SendRawTransaction accepts rlp-encode of signed transaction and sends it via RPC call
func (c *EVMClient) SendRawTransaction(ctx context.Context, tx []byte) error {
	return c.rpClient.CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(tx))
}

func (c *EVMClient) CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := c.rpClient.CallContext(ctx, &hex, "eth_call", callArgs, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (c *EVMClient) PendingCallContract(ctx context.Context, callArgs map[string]interface{}) ([]byte, error) {
	var hex hexutil.Bytes
	err := c.rpClient.CallContext(ctx, &hex, "eth_call", callArgs, "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (c *EVMClient) SignAndSendTransaction(ctx context.Context, tx CommonTransaction) (common.Hash, error) {
	id, err := c.ChainID(ctx)
	if err != nil {
		panic(err)
	}
	rawTX, err := tx.RawWithSignature(c.config.kp.PrivateKey(), id)
	if err != nil {
		return common.Hash{}, err
	}

	err = c.SendRawTransaction(ctx, rawTX)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (c *EVMClient) RelayerAddress() common.Address {
	return c.config.relayerAddress
}

func (c *EVMClient) LockNonce() {
	c.nonceLock.Lock()
}

func (c *EVMClient) UnlockNonce() {
	c.nonceLock.Unlock()
}

func (c *EVMClient) UnsafeNonce() (*big.Int, error) {
	var err error
	for i := 0; i <= 10; i++ {
		if c.nonce == nil {
			nonce, err := c.PendingNonceAt(context.Background(), c.config.kp.CommonAddress())
			if err != nil {
				time.Sleep(1)
				continue
			}
			return big.NewInt(0).SetUint64(nonce), nil
		}
		return c.nonce, nil
	}
	return nil, err
}

func (c *EVMClient) UnsafeIncreaseNonce() error {
	nonce, err := c.UnsafeNonce()
	if err != nil {
		return err
	}
	c.nonce = nonce.And(nonce, big.NewInt(1))
	return nil
}

func (c *EVMClient) GasPrice() (*big.Int, error) {
	gasPrice, err := c.SafeEstimateGas(context.TODO())
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (c *EVMClient) SafeEstimateGas(ctx context.Context) (*big.Int, error) {
	suggestedGasPrice, err := c.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	gasPrice := multiplyGasPrice(suggestedGasPrice, c.config.gasMultiplier)

	// Check we aren't exceeding our limit

	if gasPrice.Cmp(c.config.maxGasPrice) == 1 {
		return c.config.maxGasPrice, nil
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

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract common.Address, sig string, startBlock *big.Int, endBlock *big.Int) ethereum.FilterQuery {
	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []common.Address{contract},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte(sig))},
		},
	}
	return query
}
