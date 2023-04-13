package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ChainSafe/chainbridge-core/evaluate"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to the Fantom testnet
	client, err := ethclient.Dial("wss://fantom-testnet.blastapi.io/15dfa89b-e1f6-41ec-b329-dd34f509176b")
	if err != nil {
		log.Fatal(err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to retrieve latest block header: %v", err)
	}

	latestBlockNumber := header.Number

	fmt.Printf("Latest block number: %v\n", latestBlockNumber)

	// Contract address and ABI
	contractAddress := common.HexToAddress(SourceBridgeAddress)
	contractAbi, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		log.Fatal(err)
	}

	// Starting and ending block numbers
	// startBlock := uint64(14837157)
	// endBlock := uint64(14837160)

	// Get all logs for the specified contract address and block range
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: latestBlockNumber,
		// ToBlock:   new(big.Int).SetUint64(endBlock),
		Topics: [][]common.Hash{{common.HexToHash("0x968626a768e76ba1363efe44e322a6c4900c5f084e0b45f35e294dfddaa9e0d5")}},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			// fmt.Println(vLog) // pointer to event log
			// fmt.Printf("Block number: %d\n", vLog.BlockNumber)
			// fmt.Println("vLog", vLog)
			event := struct {
				OriginDomainID uint8
				DepositNonce   uint64
				Status         uint8
				DataHash       [32]byte
			}{}
			err := contractAbi.UnpackIntoInterface(&event, "ProposalEvent", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			if event.Status == 3 {
				// data := event.DataHash[:]
				// originDomainID := data[0]
				// depositNonce := binary.BigEndian.Uint64(data[1:9])
				// status := data[9]
				// dataHash := data[10:]

				// fmt.Printf("originDomainID: %d\ndepositNonce: %d\nstatus: %d\ndataHash: %x\n", originDomainID, depositNonce, status, dataHash)

				start := time.Now()
				evaluate.SetT3(event.DepositNonce, vLog.TxHash.Hex(), start)
				fmt.Printf("ProposalEvent executed: originDomainID=%v, depositNonce=%v, status=%v, dataHash=%v\n",
					event.OriginDomainID, event.DepositNonce, event.Status, hexutil.Encode(event.DataHash[:]))
			}

		default:
			header, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatalf("Failed to retrieve latest block header: %v", err)
			}

			fmt.Printf("Latest block number: %v\n", header.Number)
			time.Sleep(5 * time.Second)
		}
	}
}
