// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/relayer"
	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

const (
	DepositSignature string = "Deposit(uint8,bytes32,uint64)"
)

type EventHandler func(sourceID, destID uint8, nonce uint64, handlerContractAddress string, backend ChainReader) (*relayer.Message, error)
type EventHandlers map[ethcommon.Address]EventHandler

var BlockRetryInterval = time.Second * 5
var BlockDelay = big.NewInt(10) //TODO: move to config

type ChainReader interface {
	goeth.ChainReader
	bind.ContractFilterer
	bind.ContractCaller
	MatchResourceIDToHandlerAddress(bridgeAddress string, rID [32]byte) (string, error)
}

type EVMListener struct {
	chainReader   ChainReader
	eventHandlers EventHandlers
}

func NewEVMListener(chainReader ChainReader) *EVMListener {
	return &EVMListener{chainReader: chainReader, eventHandlers: make(map[ethcommon.Address]EventHandler)}
}

// TODO maybe it could be private
func (l *EVMListener) MatchAddressWithHandlerFunc(addr string) (EventHandler, error) {
	h, ok := l.eventHandlers[ethcommon.HexToAddress(addr)]
	if !ok {
		return nil, errors.New("no corresponding handler for this address exists")
	}
	return h, nil
}

func (l *EVMListener) RegisterHandler(address string, handler EventHandler) {
	if l.eventHandlers == nil {
		l.eventHandlers = make(map[ethcommon.Address]EventHandler)
	}
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

func (l *EVMListener) ListenToEvents(startBlock *big.Int, chainID uint8, bridgeContractAddress string, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan *relayer.Message {
	bridgeAddress := ethcommon.HexToAddress(bridgeContractAddress)
	// TODO: This channel should be closed somewhere!
	ch := make(chan *relayer.Message)
	go func() {
		for {
			select {
			case <-stopChn:
				return
			default:
				head, err := l.chainReader.HeaderByNumber(context.Background(), nil)
				if err != nil {
					log.Error().Err(err).Msg("Unable to get latest block")
					time.Sleep(BlockRetryInterval)
					continue
				}
				// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
				if big.NewInt(0).Sub(head.Number, startBlock).Cmp(BlockDelay) == -1 {
					time.Sleep(BlockRetryInterval)
					continue
				}
				query := buildQuery(bridgeAddress, DepositSignature, startBlock, startBlock)
				logs, err := l.chainReader.FilterLogs(context.Background(), query)
				if err != nil {
					// Filtering logs error really can appear only on wrong configuration or temporary network problem
					// so i do no see any reason to break execution
					log.Error().Err(err).Str("ChainID", string(chainID)).Msgf("Unable to filter logs")
					continue
				}
				for _, eventLog := range logs {
					destId := uint8(eventLog.Topics[1].Big().Uint64())
					rId := eventLog.Topics[2]
					nonce := eventLog.Topics[3].Big().Uint64()

					addr, err := l.chainReader.MatchResourceIDToHandlerAddress(bridgeContractAddress, rId)
					if err != nil {
						errChn <- err
						log.Error().Err(err)
						return
					}

					eventHandler, err := l.MatchAddressWithHandlerFunc(addr)
					if err != nil {
						errChn <- err
						log.Error().Err(err).Msgf("failed to get handler from resource ID %x, reason: %w", rId, err)
						return
					}

					m, err := eventHandler(chainID, destId, nonce, addr, l.chainReader)
					if err != nil {
						errChn <- err
						log.Error().Err(err)
						return
					}
					log.Debug().Msgf("Resolved message %+v in block %s", m, startBlock.String())
					ch <- m
				}

				if startBlock.Int64()%20 == 0 {
					// Logging process every 20 bocks to exclude spam
					log.Debug().Str("block", startBlock.String()).Uint8("chainID", chainID).Msg("Queried block for deposit events")
				}
				// TODO: We can store blocks to DB inside listener or make listener send something to channel each block to save it.
				//Write to block store. Not a critical operation, no need to retry
				err = blockstore.StoreBlock(kvrw, startBlock, chainID)
				if err != nil {
					log.Error().Str("block", startBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
				}

				// Goto next block
				startBlock.Add(startBlock, big.NewInt(1))
			}
		}
	}()
	return ch
}
