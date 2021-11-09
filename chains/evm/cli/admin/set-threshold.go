package admin

import (
	"fmt"
	"math/big"

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

var setThresholdCmd = &cobra.Command{
	Use:   "set-threshold",
	Short: "Set a new relayer vote threshold",
	Long:  "Set a new relayer vote threshold",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return SetThresholdCMD(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
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

func SetThresholdCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	log.Debug().Msgf(`
Setting new threshold
Threshold: %d
Bridge address: %s`, RelayerThreshold, Bridge)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})
	setThresholdInput, err := calls.PrepareSetThresholdInput(big.NewInt(0).SetUint64(RelayerThreshold))
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare set threshold input error: %v", err))
		return err
	}
	_, err = calls.Transact(ethClient, txFabric, gasPricer, &bridgeAddr, setThresholdInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Error().Err(fmt.Errorf("transact error: %v", err))
		return err
	}
	log.Info().Msgf("New threshold set to %v", RelayerThreshold)
	return nil
}
