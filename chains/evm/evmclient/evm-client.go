package evmclient

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
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
	rpClient *rpc.Client
}

type CommonTransaction interface {
	// Hash returns the transaction hash.
	Hash() common.Hash
	// RawWithSignature mostly copies WithSignature interface of type.Transaction from go-ethereum,
	// but return raw rlp encoded signed transaction to be compatible and interchangeable between different go-ethereum implementations
	// WithSignature returns a new transaction with the given signature.
	// This signature needs to be in the [R || S || V] format where V is 0 or 1.
	RawWithSignature(types.Signer, []byte) ([]byte, error)
}

func NewEVMClient(endpoint string) (*EVMClient, error) {
	log.Info().Str("url", endpoint).Msg("Connecting to evm chain...")
	rpcClient, err := rpc.DialContext(context.TODO(), endpoint)
	if err != nil {
		return nil, err
	}
	return &EVMClient{
		Client:   ethclient.NewClient(rpcClient),
		rpClient: rpcClient,
	}, nil
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

func (c *EVMClient) FetchDepositLogs(ctx context.Context, contractAddress string, startBlock *big.Int, endBlock *big.Int) ([]*listener.DepositLogs, error) {
	logs, err := c.FilterLogs(ctx, buildQuery(common.HexToAddress(contractAddress), DepositSignature, startBlock, endBlock))
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
	h := tx.Hash()
	sig, err := crypto.Sign(h[:], prvKey)
	if err != nil {
		return common.Hash{}, err
	}

	rawTX, err := tx.RawWithSignature(types.HomesteadSigner{}, sig)
	if err != nil {
		return common.Hash{}, err
	}

	err = c.SendRawTransaction(ctx, rawTX)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
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
