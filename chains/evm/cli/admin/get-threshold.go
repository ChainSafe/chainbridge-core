package admin

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/contracts"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var getThresholdCmd = &cobra.Command{
	Use:   "get-threshold",
	Short: "get relayer vote threshold",
	Long:  "get relayer vote threshold",
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
		return GetThresholdCMD(cmd, args, bridgeContract)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateGetThresholdFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessGetThresholdFlags(cmd, args)
		return nil
	},
}

func BindGetThresholdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "bridge")
}
func init() {
	BindGetThresholdFlags(getThresholdCmd)
}

func ValidateGetThresholdFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func ProcessGetThresholdFlags(cmd *cobra.Command, args []string) {
	bridgeAddr = common.HexToAddress(Bridge)
}

func GetThresholdCMD(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
getting threshold
Bridge address: %s`, Bridge)
	threshold, err := contract.GetThreshold()
	if err != nil {
		log.Error().Err(fmt.Errorf("transact error: %v", err))
		return err
	}
	log.Info().Msgf("Relayer threshold for the bridge %v is %v", Bridge, threshold)
	return nil
}
