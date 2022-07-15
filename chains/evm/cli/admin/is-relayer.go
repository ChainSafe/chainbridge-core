package admin

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var isRelayerCmd = &cobra.Command{
	Use:   "is-relayer",
	Short: "Check if an address is registered as a relayer",
	Long:  "The is-relayer subcommand checks if an address is registered as a relayer",
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
		return IsRelayer(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateIsRelayerFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessIsRelayerFlags(cmd, args)
		return nil
	},
}

func BindIsRelayerFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Relayer, "relayer", "", "Address to check")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "relayer", "bridge")
}

func init() {
	BindIsRelayerFlags(isRelayerCmd)
}

func ValidateIsRelayerFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Relayer) {
		return fmt.Errorf("invalid relayer address %s", Relayer)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func ProcessIsRelayerFlags(cmd *cobra.Command, args []string) {
	RelayerAddr = common.HexToAddress(Relayer)
	BridgeAddr = common.HexToAddress(Bridge)
}

func IsRelayer(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
	Checking relayer
	Relayer address: %s
	Bridge address: %s`, Relayer, Bridge)

	isRelayer, err := contract.IsRelayer(RelayerAddr)
	if err != nil {
		return err
	}

	if !isRelayer {
		log.Info().Msgf("Address %s is NOT relayer", RelayerAddr.String())
	} else {
		log.Info().Msgf("Address %s is relayer", RelayerAddr.String())
	}
	return nil
}
