package evmclient

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

type EVMClient struct {
	*ethclient.Client
	rpClient  *rpc.Client
	config    *EVMConfig
	nonce     *big.Int
	nonceLock sync.Mutex
	gasPrice  *big.Int
}

type CommonTransaction interface {
	// Hash returns the transaction hash.
	Hash() common.Hash

	// RawWithSignature Returns signed transaction by provided private key
	RawWithSignature(key *ecdsa.PrivateKey, chainID *big.Int) ([]byte, error)
}

func NewEVMClient() *EVMClient {
	return &EVMClient{}
}

func NewEVMClientFromParams(url string, privateKey *ecdsa.PrivateKey, gasPrice *big.Int) (*EVMClient, error) {
	rpcClient, err := rpc.DialContext(context.TODO(), url)
	if err != nil {
		return nil, err
	}
	kp := secp256k1.NewKeypair(*privateKey)
	c := &EVMClient{}
	c.Client = ethclient.NewClient(rpcClient)
	c.rpClient = rpcClient
	c.config = &EVMConfig{}
	c.config.kp = kp
	c.gasPrice = gasPrice
	return c, nil
}

func (c *EVMClient) Configurate(path string, name string) error {
	rawCfg, err := GetConfig(path, name)
	if err != nil {
		return err
	}
	cfg, err := ParseConfig(rawCfg)
	if err != nil {
		return err
	}
	c.config = cfg
	generalConfig := cfg.SharedEVMConfig.GeneralChainConfig

	kp, err := keystore.KeypairFromAddress(generalConfig.From, keystore.EthChain, generalConfig.KeystorePath, generalConfig.Insecure)
	if err != nil {
		panic(err)
	}
	krp := kp.(*secp256k1.Keypair)
	c.config.kp = krp

	log.Info().Str("url", generalConfig.Endpoint).Msg("Connecting to evm chain...")
	rpcClient, err := rpc.DialContext(context.TODO(), generalConfig.Endpoint)
	if err != nil {
		return err
	}
	c.Client = ethclient.NewClient(rpcClient)
	c.rpClient = rpcClient

	if generalConfig.LatestBlock {
		curr, err := c.LatestBlock()
		if err != nil {
			return err
		}
		cfg.SharedEVMConfig.StartBlock = curr
	}

	return nil

}

type headerNumber struct {
	Number *big.Int `json:"number"           gencodec:"required"`
}

func (h *headerNumber) UnmarshalJSON(input []byte) error {
	type headerNumber struct {
		Number *hexutil.Big `json:"number" gencodec:"required"`
	}
	var dec headerNumber
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	h.Number = (*big.Int)(dec.Number)
	return nil
}

// LatestBlock returns the latest block from the current chain
func (c *EVMClient) LatestBlock() (*big.Int, error) {
	var head *headerNumber
	err := c.rpClient.CallContext(context.Background(), &head, "eth_getBlockByNumber", toBlockNumArg(nil), false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	if err != nil {
		return nil, err
	}
	return head.Number, nil
}

func (c *EVMClient) WaitAndReturnTxReceipt(h common.Hash) (*types.Receipt, error) {
	retry := 50
	for retry > 0 {
		receipt, err := c.Client.TransactionReceipt(context.Background(), h)
		if err != nil {
			log.Error().Err(err).Msgf("error getting tx receipt %s", h.String())
			retry--
			time.Sleep(5 * time.Second)
			continue
		}
		if receipt.Status != 1 {
			return receipt, errors.New("transaction failed on chain")
		}
		return receipt, nil
	}
	return nil, errors.New("tx did not appear")
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

func (c *EVMClient) FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error) {
	return c.FilterLogs(ctx, buildQuery(contractAddress, event, startBlock, endBlock))
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

func (c *EVMClient) CallContext(ctx context.Context, target interface{}, rpcMethod string, args ...interface{}) error {
	err := c.rpClient.CallContext(ctx, target, rpcMethod, args)
	if err != nil {
		return err
	}
	return nil
}

func (c *EVMClient) PendingCallContract(ctx context.Context, callArgs map[string]interface{}) ([]byte, error) {
	var hex hexutil.Bytes
	err := c.rpClient.CallContext(ctx, &hex, "eth_call", callArgs, "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (c *EVMClient) From() common.Address {
	return c.config.kp.CommonAddress()
}

func (c *EVMClient) SignAndSendTransaction(ctx context.Context, tx CommonTransaction) (common.Hash, error) {
	id, err := c.ChainID(ctx)
	if err != nil {
		//panic(err)
		// Probably chain does not support ChainID eg. CELO
		id = nil
	}
	rawTx, err := tx.RawWithSignature(c.config.kp.PrivateKey(), id)
	if err != nil {
		return common.Hash{}, err
	}

	err = c.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (c *EVMClient) RelayerAddress() common.Address {
	return c.config.kp.CommonAddress()
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
			c.nonce = big.NewInt(0).SetUint64(nonce)
			return c.nonce, nil
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
	c.nonce = nonce.Add(nonce, big.NewInt(1))
	return nil
}

func (c *EVMClient) BaseFee() (*big.Int, error) {
	head, err := c.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return head.BaseFee, nil
}

func (c *EVMClient) GasPrice() (*big.Int, error) {
	// NOTE: do we want this? Is this gas price method meant to have the gas price stay the same once it is set?
	if c.gasPrice != nil {
		return c.gasPrice, nil
	}
	gasPrice, err := c.SafeEstimateGas(context.TODO())
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (c *EVMClient) EstimateGasLondon(ctx context.Context, baseFee *big.Int) (*big.Int, *big.Int, error) {
	var maxPriorityFeePerGas *big.Int
	var maxFeePerGas *big.Int

	var sharedEVMConfig config.SharedEVMConfig = c.config.SharedEVMConfig
	if sharedEVMConfig.MaxGasPrice.Cmp(baseFee) < 0 {
		maxPriorityFeePerGas = big.NewInt(1)
		maxFeePerGas = new(big.Int).Add(sharedEVMConfig.MaxGasPrice, maxPriorityFeePerGas)
		return maxPriorityFeePerGas, maxFeePerGas, nil
	}

	maxPriorityFeePerGas, err := c.SuggestGasTipCap(context.TODO())
	if err != nil {
		return nil, nil, err
	}

	maxFeePerGas = new(big.Int).Add(
		maxPriorityFeePerGas,
		new(big.Int).Mul(baseFee, big.NewInt(2)),
	)

	if maxFeePerGas.Cmp(maxPriorityFeePerGas) < 0 {
		return nil, nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", maxFeePerGas, maxPriorityFeePerGas)
	}

	// Check we aren't exceeding our limit
	if maxFeePerGas.Cmp(sharedEVMConfig.MaxGasPrice) == 1 {
		maxPriorityFeePerGas.Sub(sharedEVMConfig.MaxGasPrice, baseFee)
		maxFeePerGas = sharedEVMConfig.MaxGasPrice
	}
	return maxPriorityFeePerGas, maxFeePerGas, nil
}

func (c *EVMClient) SafeEstimateGas(ctx context.Context) (*big.Int, error) {
	suggestedGasPrice, err := c.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("Suggested GP %s", suggestedGasPrice.String())
	var gasPrice *big.Int
	if c.config.SharedEVMConfig.GasMultiplier != nil {
		gasPrice = multiplyGasPrice(suggestedGasPrice, c.config.SharedEVMConfig.GasMultiplier)
	}
	// Check we aren't exceeding our limit
	if c.config.SharedEVMConfig.MaxGasPrice != nil {
		if gasPrice.Cmp(c.config.SharedEVMConfig.MaxGasPrice) == 1 {
			return c.config.SharedEVMConfig.MaxGasPrice, nil
		}
	}
	return gasPrice, nil
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

func (c *EVMClient) GetConfig() *EVMConfig {
	return c.config
}
