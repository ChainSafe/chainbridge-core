package admin

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause deposits and proposals",
	Long:  "The pause subcommand pauses deposits and proposals",
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
		return PauseCmd(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidatePauseFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessPauseFlags(cmd, args)

		return nil
	},
}

func BindPauseFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "bridge")
}

func init() {
	BindPauseFlags(pauseCmd)
}

func ValidatePauseFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address: %s", Bridge)
	}
	return nil
}

func ProcessPauseFlags(cmd *cobra.Command, args []string) {
	BridgeAddr = common.HexToAddress(Bridge)
}

func PauseCmd(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	hash, err := contract.Pause(transactor.TransactOptions{GasLimit: gasLimit})
	if err != nil {
		log.Error().Err(fmt.Errorf("admin pause error: %v", err))
		return err
	}

	log.Info().Msgf("successfully paused bridge: %s; tx hash: %s", Bridge, hash.Hex())
	return nil
}
