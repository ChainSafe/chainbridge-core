package utils

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var hashListCmd = &cobra.Command{
	Use:   "hash-list",
	Short: "List tx hashes within N number of blocks",
	Long:  "The hash-list subcommand accepts a starting block to query, loops over N number of blocks past it, then prints this list of blocks to review hashes contained within",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return HashListCmd(cmd, args)
	},
}

func BindHashListCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&BlockNumber, "blockNumber", "", "Block number to start at")
	cmd.Flags().StringVar(&NumberOfBlocks, "numberOfBlocks", "", "Number of blocks past the provided blockNumber to review")
	flags.MarkFlagsAsRequired(cmd, "blockNumber", "numberOfBlocks")
}

func init() {
	BindHashListCmdFlags(hashListCmd)
}

func HashListCmd(cmd *cobra.Command, args []string) error {
	// fetch global flag values
	url, _, _, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	// convert NumberOfBlocks string to int for looping
	numberOfBlocks, err := strconv.Atoi(NumberOfBlocks)
	if err != nil {
		log.Error().Err(fmt.Errorf("error converting NumberOfBlocks string -> int: %v", err))
		return err
	}

	// convert block number to string
	blockNumberBigInt, _ := new(big.Int).SetString(BlockNumber, 10)

	// declare empty slice of blocks to hold blocks for printing all at once
	blockSlice := make([]*types.Block, 0)

	// loop over blocks provided by user
	// check block by hash
	// see if transaction block data is there
	for i := 0; i < numberOfBlocks; i++ {
		log.Debug().Msgf("blockNum: %v", blockNumberBigInt)

		// convert string block number to big.Int
		blockNumberBigInt.Add(blockNumberBigInt, big.NewInt(1))

		block, err := ethClient.BlockByNumber(context.Background(), blockNumberBigInt)
		if err != nil {
			log.Error().Err(fmt.Errorf("block by hash error: %v", err))

			// will return early and not print debug log if block not found
			// Error: not found

			// return err
		}

		// performance: append to block to slice of blocks to return all at once
		// rather than printing each, one-by-one
		blockSlice = append(blockSlice, block)
	}

	// log full struct of block
	log.Debug().Msgf("block slice: %+v", blockSlice)

	return nil
}
