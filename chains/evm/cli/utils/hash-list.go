package utils

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
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

func BindHashListFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&BlockNumber, "block-number", "", "Block number to start at")
	cmd.Flags().StringVar(&Blocks, "blocks", "", "Number of blocks past the provided block-number to review")
	flags.MarkFlagsAsRequired(cmd, "block-number", "blocks")
}

func init() {
	BindHashListFlags(hashListCmd)
}

func HashListCmd(cmd *cobra.Command, args []string) error {
	// fetch global flag values
	url, _, _, senderKeyPair, _, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	// convert Blocks string to int for looping
	numBlocks, err := strconv.Atoi(Blocks)
	if err != nil {
		log.Error().Err(fmt.Errorf("error converting NumberOfBlocks string -> int: %v", err))
		return err
	}

	// convert block number to string
	blockNumberBigInt, _ := new(big.Int).SetString(BlockNumber, 10)

	// loop over blocks provided by user
	// check block by hash
	// see if transaction block data is there
	for i := 0; i < numBlocks; i++ {
		log.Debug().Msgf("Block Number: %v", blockNumberBigInt)

		// convert string block number to big.Int
		blockNumberBigInt.Add(blockNumberBigInt, big.NewInt(1))

		block, err := ethClient.BlockByNumber(context.Background(), blockNumberBigInt)
		if err != nil {
			log.Error().Err(fmt.Errorf("block by hash error: %v", err))

			// will return early and not print debug log if block not found
			// Error: not found

			return err
		}

		// loop over all transactions within block
		// add newline for readability
		for _, tx := range block.Body().Transactions {
			log.Debug().Msgf("Tx hashes: %v\n", tx.Hash())
		}
	}
	return nil
}
