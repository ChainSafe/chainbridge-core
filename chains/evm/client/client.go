package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ContractCaller interface {
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
}

type TransactionDispatcher interface {
	WaitAndReturnTxReceipt(h common.Hash) (*types.Receipt, error)
	SignAndSendTransaction(ctx context.Context, tx CommonTransaction) (common.Hash, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	GetTransactionByHash(h common.Hash) (tx *types.Transaction, isPending bool, err error)
	UnsafeNonce() (*big.Int, error)
	LockNonce()
	UnlockNonce()
	UnsafeIncreaseNonce() error
	From() common.Address
}

type Client interface {
	ContractCaller
	TransactionDispatcher
}

type EVMClient struct {
	*ethclient.Client
	signer     Signer
	gethClient *gethclient.Client
	rpClient   *rpc.Client
	nonce      *big.Int
	nonceLock  sync.Mutex
}

type Signer interface {
	CommonAddress() common.Address

	// Sign calculates an ECDSA signature.
	// The produced signature must be in the [R || S || V] format where V is 0 or 1.
	Sign(digestHash []byte) ([]byte, error)
}

type CommonTransaction interface {
	// Hash returns the transaction hash.
	Hash() common.Hash

	// RawWithSignature Returns signed transaction by provided signer
	RawWithSignature(signer Signer, domainID *big.Int) ([]byte, error)
}

// NewEVMClient creates a client for EVMChain with provided signer
func NewEVMClient(url string, signer Signer) (*EVMClient, error) {
	rpcClient, err := rpc.DialContext(context.TODO(), url)
	if err != nil {
		return nil, err
	}
	c := &EVMClient{}
	c.Client = ethclient.NewClient(rpcClient)
	c.gethClient = gethclient.New(rpcClient)
	c.rpClient = rpcClient
	c.signer = signer
	return c, nil
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

func (c *EVMClient) GetTransactionByHash(h common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return c.Client.TransactionByHash(context.Background(), h)
}

func (c *EVMClient) FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error) {
	logs, err := c.FilterLogs(ctx, buildQuery(contractAddress, event, startBlock, endBlock))
	if err != nil {
		return []types.Log{}, err
	}

	validLogs := make([]types.Log, 0)
	for _, log := range logs {
		if log.Removed {
			continue
		}

		validLogs = append(validLogs, log)
	}
	return validLogs, nil
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
	err := c.rpClient.CallContext(ctx, target, rpcMethod, args...)
	if err != nil {
		return err
	}
	return nil
}

func (c *EVMClient) From() common.Address {
	return c.signer.CommonAddress()
}

func (c *EVMClient) SignAndSendTransaction(ctx context.Context, tx CommonTransaction) (common.Hash, error) {
	id, err := c.ChainID(ctx)
	if err != nil {
		// panic(err)
		// Probably chain does not support chainID eg. CELO
		id = nil
	}
	rawTx, err := tx.RawWithSignature(c.signer, id)
	if err != nil {
		return common.Hash{}, err
	}
	err = c.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (c *EVMClient) WaitAndReturnTxReceipt(h common.Hash) (*types.Receipt, error) {
	retry := 50
	for retry > 0 {
		receipt, err := c.Client.TransactionReceipt(context.Background(), h)
		if err != nil {
			retry--
			time.Sleep(5 * time.Second)
			continue
		}
		if receipt.Status != 1 {
			return receipt, fmt.Errorf("transaction failed on chain. Receipt status %v", receipt.Status)
		}
		return receipt, nil
	}
	return nil, errors.New("tx did not appear")
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
			nonce, err := c.PendingNonceAt(context.Background(), c.signer.CommonAddress())
			if err != nil {
				time.Sleep(1 * time.Second)
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
