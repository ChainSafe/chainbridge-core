package evmd

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/Bridge"
	"github.com/ChainSafe/chainbridgev2/relayer"
	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	Deposit string = "Deposit(uint8,bytes32,uint64)"
)

type Handler func(sourceID, destID uint8, nonce uint64, handlerContractAddress common.Address, backend bind.ContractBackend) (relayer.XCMessager, error)
type Handlers map[ethcommon.Address]Handler

type IBridge interface {
	ResourceIDToHandlerAddress(opts *bind.CallOpts, arg0 [32]byte) (ethcommon.Address, error)
}

var ExpectedBlockTime = time.Second

func NewListener(client *Client, bridgeAddress ethcommon.Address, chainID uint8) *Listener {
	bridgeContract, err := Bridge.NewBridge(bridgeAddress, client)
	if err != nil {
		panic(err)
	}

	return &Listener{bridgeContract: bridgeContract, BridgeContractAddress: bridgeAddress, client: client, chainID: chainID}
}

type Listener struct {
	stop                  <-chan struct{}
	sysErr                chan<- error // Reports fatal error to core
	bridgeContract        IBridge
	BridgeContractAddress ethcommon.Address
	handlers              Handlers
	currentBlock          *big.Int
	client                *Client
	chainID               uint8
}

func (l *Listener) LatestBlock() (*big.Int, error) {
	header, err := l.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}
func (l *Listener) GetBridgeAddress() ethcommon.Address {
	return l.BridgeContractAddress
}
func (l *Listener) StoreCurrentBlock(lb *big.Int) error {
	l.currentBlock = lb
	return nil
}

func (l *Listener) MatchResourceIDToHandlerAddress(rID [32]byte) (ethcommon.Address, error) {
	return l.bridgeContract.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
}

func (l *Listener) MatchAddressWithHandlerFunc(addr ethcommon.Address) (Handler, error) {
	h, ok := l.handlers[addr]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return h, nil
}

func (l *Listener) RegisterHandler(address ethcommon.Address, handler Handler) {
	l.handlers[address] = handler
}

func (l *Listener) LogsForBlock(ctx context.Context, latestBlock *big.Int) ([]types.Log, error) {
	Deposit := "Deposit(uint8,bytes32,uint64)"
	query := buildQuery(l.GetBridgeAddress(), Deposit, latestBlock, latestBlock)
	return l.client.FilterLogs(ctx, query)
}

func (l *Listener) GetContractBackend() bind.ContractBackend {
	return l.client
}

func (l *Listener) GetChainID() uint8 {
	return l.chainID
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig string, startBlock *big.Int, endBlock *big.Int) goeth.FilterQuery {
	query := goeth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []ethcommon.Address{contract},
		Topics: [][]ethcommon.Hash{
			{crypto.Keccak256Hash([]byte(sig))},
		},
	}
	return query
}
