package utils

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate transaction invocation",
	Long:  "Replay a failed transaction by simulating invocation; not state-altering",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return SimulateCmd(cmd)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateSimulateFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessSimulateFlags(cmd, args)
		return nil
	},
}

func BindSimulateCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&TxHash, "txHash", "", "transaction hash")
	cmd.Flags().StringVar(&BlockNumber, "blockNumber", "", "block number")
	cmd.Flags().StringVar(&FromAddress, "fromAddress", "", "address of sender")
	flags.MarkFlagsAsRequired(cmd, "txHash", "blockNumber", "fromAddress")
}

func init() {
	BindSimulateCmdFlags(simulateCmd)
}

func ValidateSimulateFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(TxHash) {
		return fmt.Errorf("invalid tx hash %s", TxHash)
	}
	if !common.IsHexAddress(FromAddress) {
		return fmt.Errorf("invalid from address %s", FromAddress)
	}
	return nil
}

func ProcessSimulateFlags(cmd *cobra.Command, args []string) {
	txHash = common.HexToHash(TxHash)
	fromAddr = common.HexToAddress(FromAddress)
}

func SimulateCmd(cmd *cobra.Command) error {
	// fetch global flag values
	url, _, _, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	// convert string block number to big.Int
	blockNumberBigInt, _ := new(big.Int).SetString(BlockNumber, 10)

	log.Debug().Msgf(`
Simulating transaction
Tx hash: %s
Block number: %v
From address: %s`,
		TxHash, blockNumberBigInt, FromAddress)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	data, err := calls.Simulate(ethClient, blockNumberBigInt, txHash, fromAddr)
	if err != nil {
		log.Error().Err(fmt.Errorf("[utils] simulate transact error: %v", err))
		return err
	}

	log.Debug().Msgf("data: %s", string(data))

	return nil
}
