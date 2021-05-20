// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package evm

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type EventListener interface {
	ListenToEvents(startBlock *big.Int, chainID uint8, bridgeContractAddress string, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan *relayer.Message
}

type ProposalVoter interface {
	VoteProposal(message *relayer.Message, bridgeAddress string) error
}

// EVMChain is struct that aggregates all data required for
type EVMChain struct {
	listener              EventListener // Rename
	writer                ProposalVoter
	chainID               uint8
	kvdb                  blockstore.KeyValueReaderWriter
	bridgeContractAddress string
}

func NewEVMChain(dr EventListener, writer ProposalVoter, kvdb blockstore.KeyValueReaderWriter, bridgeContractAddress string, chainID uint8) *EVMChain {
	return &EVMChain{listener: dr, writer: writer, kvdb: kvdb, bridgeContractAddress: bridgeContractAddress, chainID: chainID}
}

// PollEvents is the goroutine that polling blocks and searching Deposit Events in them. Event then sent to eventsChan
func (c *EVMChain) PollEvents(stop <-chan struct{}, sysErr chan<- error, eventsChan chan *relayer.Message) {
	log.Info().Msg("Polling Blocks...")
	// Handler chain specific configs and flags
	b, err := blockstore.GetLastStoredBlock(c.kvdb, c.chainID)
	if err != nil {
		sysErr <- fmt.Errorf("error %w on getting last stored block", err)
		return
	}
	ech := c.listener.ListenToEvents(b, c.chainID, c.bridgeContractAddress, c.kvdb, stop, sysErr)
	for {
		select {
		case <-stop:
			return
		case newEvent := <-ech:
			// Here we can place middlewares for custom logic?
			eventsChan <- newEvent
			continue
		}
	}
}

// Write function pass XCMessager to underlying chain writer
func (c *EVMChain) Write(msg *relayer.Message) error {
	return c.writer.VoteProposal(msg, c.bridgeContractAddress)
}

func (c *EVMChain) ChainID() uint8 {
	return c.chainID
}

// TODO: should be moved somewhere else
type Proposal struct {
	Source         uint8  // Source where message was initiated
	Destination    uint8  // Destination chain of message
	DepositNonce   uint64 // Nonce for the deposit
	ResourceId     [32]byte
	Payload        []interface{} // data associated with event sequence
	Data           []byte
	DataHash       common.Hash
	HandlerAddress common.Address
}

func GetIDAndNonce(p *Proposal) *big.Int {
	data := bytes.Buffer{}
	bn := big.NewInt(0).SetUint64(p.DepositNonce).Bytes()
	data.Write(bn)
	data.Write([]byte{p.Source})
	return big.NewInt(0).SetBytes(data.Bytes())
}
