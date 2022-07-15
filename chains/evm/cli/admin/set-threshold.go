package admin

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setThresholdCmd = &cobra.Command{
	Use:   "set-threshold",
	Short: "Set a new relayer vote threshold",
	Long:  "The set-threshold subcommand sets a new relayer vote threshold",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c, prepare)
		if err != nil {
			return err
		}
		return SetThresholdCMD(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
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
	cmd.Flags().Uint64Var(&RelayerThreshold, "threshold", 0, "New relayer threshold")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
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
	BridgeAddr = common.HexToAddress(Bridge)
}

func SetThresholdCMD(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
Setting new threshold
Threshold: %d
Bridge address: %s`, RelayerThreshold, Bridge)
	_, err := contract.AdminChangeRelayerThreshold(RelayerThreshold, transactor.TransactOptions{GasLimit: gasLimit})
	if err != nil {
		return err
	}
	return nil
}
