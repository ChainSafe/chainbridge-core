package utils

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate the invocation of a transaction",
	Long:  "The simulate subcommand simulates a transaction result by simulating invocation; not state-altering",
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

func BindSimulateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&TxHash, "tx-hash", "", "Transaction hash")
	cmd.Flags().StringVar(&BlockNumber, "block-number", "", "Block number")
	cmd.Flags().StringVar(&FromAddress, "from", "", "Address of sender")
	flags.MarkFlagsAsRequired(cmd, "tx-hash", "block-number", "from")
}

func init() {
	BindSimulateFlags(simulateCmd)
}

func ValidateSimulateFlags(cmd *cobra.Command, args []string) error {
	_, err := hexutil.Decode(TxHash)
	if err != nil {
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
	url, _, _, senderKeyPair, _, err := flags.GlobalFlagValues(cmd)
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
