package admin

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
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
		txFabric := evmtransaction.NewTransaction
		return GetThresholdCMD(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
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

func GetThresholdCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	log.Debug().Msgf(`
getting threshold
Bridge address: %s`, Bridge)

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

	threshold, err := calls.GetThreshold(ethClient, &bridgeAddr)
	if err != nil {
		log.Error().Err(fmt.Errorf("transact error: %v", err))
		return err
	}
	log.Info().Msgf("Relayer threshold for the bridge %v is %v", Bridge, threshold)
	return nil
}
