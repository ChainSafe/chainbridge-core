package utils

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var hashListCmd = &cobra.Command{
	Use:   "hashList",
	Short: "List tx hashes",
	Long:  "List tx hashes",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return HashListCmd(cmd, args)
	},
}

func BindHashListCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&BlockNumber, "blockNumber", "", "block number")
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

	blockNum, err := strconv.Atoi(BlockNumber)
	if err != nil {
		log.Error().Err(fmt.Errorf("block string->int conversion error: %v", err))
		return err
	}

	blockNumStr := strconv.Itoa(blockNum)
	blockNumberBigInt, _ := new(big.Int).SetString(blockNumStr, 10)

	// check block by hash
	// see if transaction block data is there
	for i := 0; i < 50; i++ {
		log.Debug().Msgf("blockNum: %v", blockNumberBigInt)

		// convert string block number to big.Int
		// ignore success bool

		blockNumberBigInt.Add(blockNumberBigInt, big.NewInt(1))

		block, err := ethClient.BlockByNumber(context.Background(), blockNumberBigInt)
		if err != nil {
			log.Error().Err(fmt.Errorf("block by hash error: %v", err))

			// will return early and not print debug log if block not found
			// Error: not found

			// return err
		}

		log.Debug().Msgf("block: %v", block)
	}
	return nil
}
