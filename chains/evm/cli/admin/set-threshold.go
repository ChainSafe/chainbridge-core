package admin

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setThresholdCmd = &cobra.Command{
	Use:   "set-threshold",
	Short: "Set a new relayer vote threshold",
	Long:  "Set a new relayer vote threshold",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		bridgeContract, err := contracts.InitializeBridgeContract(
			url, gasLimit, gasPrice, senderKeyPair, bridgeAddr,
		)
		if err != nil {
			return err
		}
		return SetThresholdCMD(cmd, args, bridgeContract)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateSetThresholdFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessSetThresholdFlags(cmd, args)
		return nil
	},
}

func BindSetThresholdFlags(cmd *cobra.Command) {
	cmd.Flags().Uint64Var(&RelayerThreshold, "threshold", 0, "new relayer threshold")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "threshold", "bridge")
}
func init() {
	BindSetThresholdFlags(setThresholdCmd)
}

func ValidateSetThresholdFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func ProcessSetThresholdFlags(cmd *cobra.Command, args []string) {
	bridgeAddr = common.HexToAddress(Bridge)
}

func SetThresholdCMD(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
Setting new threshold
Threshold: %d
Bridge address: %s`, RelayerThreshold, Bridge)
	_, err := contract.SetThresholdInput(RelayerThreshold, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	return nil
}
