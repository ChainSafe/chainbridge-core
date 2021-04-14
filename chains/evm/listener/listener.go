package listener

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridgev2/relayer"
	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

const (
	DepositSignature string = "Deposit(uint8,bytes32,uint64)"
)

type EventHandler func(sourceID, destID uint8, nonce uint64, handlerContractAddress string, backend ChainClient) (relayer.XCMessager, error)
type EventHandlers map[ethcommon.Address]EventHandler

type ChainClient interface {
	//goeth.LogFilterer
	//goeth.ChainReader
	bind.ContractCaller
	FilterLogs(ctx context.Context, q goeth.FilterQuery) ([]types.Log, error)
	MatchResourceIDToHandlerAddress(rID [32]byte, bridgeAddress string) (string, error)
}

type EVMListener struct {
	chainReader           ChainClient
	bridgeContractAddress ethcommon.Address
	eventHandlers         EventHandlers
	chainID               uint8
}

func NewEVMListener(chainReader ChainClient, bridgeContractAddress string, chainID uint8) *EVMListener {
	return &EVMListener{chainReader: chainReader, bridgeContractAddress: ethcommon.HexToAddress(bridgeContractAddress), chainID: chainID}
}

func (l *EVMListener) GetDepositEventsForBlockRange(blockFrom, blockTo *big.Int) ([]relayer.XCMessager, error) {
	query := buildQuery(l.bridgeContractAddress, DepositSignature, blockFrom, blockTo)
	logs, err := l.chainReader.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("unable to Filter Logs: %w", err)
	}
	if len(logs) == 0 {
		return []relayer.XCMessager{}, nil
	}
	msg := make([]relayer.XCMessager, 0)
	// read through the log events and handle their deposit event if handler is recognized
	for _, eventLog := range logs {
		destId := uint8(eventLog.Topics[1].Big().Uint64())
		rId := eventLog.Topics[2]
		nonce := eventLog.Topics[3].Big().Uint64()

		addr, err := l.chainReader.MatchResourceIDToHandlerAddress(rId, l.bridgeContractAddress.String())
		if err != nil {
			return nil, err
		}

		eventHandler, err := l.MatchAddressWithHandlerFunc(addr)
		if err != nil {
			return nil, fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rId, err)
		}

		m, err := eventHandler(l.chainID, destId, nonce, addr, l.chainReader)
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("Resolved message %+v", m)
		msg = append(msg, m)
	}
	return msg, nil
}

func (l *EVMListener) MatchAddressWithHandlerFunc(addr string) (EventHandler, error) {
	h, ok := l.eventHandlers[ethcommon.HexToAddress(addr)]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return h, nil
}

func (l *EVMListener) RegisterHandler(address string, handler EventHandler) {
	l.eventHandlers[ethcommon.HexToAddress(address)] = handler
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
