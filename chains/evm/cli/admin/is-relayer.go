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

var isRelayerCmd = &cobra.Command{
	Use:   "is-relayer",
	Short: "Check if an address is registered as a relayer",
	Long:  "Check if an address is registered as a relayer",
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
		return IsRelayer(cmd, args, bridgeContract)
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
	cmd.Flags().StringVar(&Relayer, "relayer", "", "address to check")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
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
	relayerAddr = common.HexToAddress(Relayer)
	bridgeAddr = common.HexToAddress(Bridge)
}

func IsRelayer(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
	Checking relayer
	Relayer address: %s
	Bridge address: %s`, Relayer, Bridge)

	isRelayer, err := contract.IsRelayer(relayerAddr)
	if err != nil {
		return err
	}

	if !isRelayer {
		log.Info().Msgf("Address %s is NOT relayer", relayerAddr.String())
	} else {
		log.Info().Msgf("Address %s is relayer", relayerAddr.String())
	}
	return nil
}
