package admin

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addRelayerCmd = &cobra.Command{
	Use:   "add-relayer",
	Short: "Add a new relayer",
	Long:  "The add-relayer subcommand sets an address as a bridge relayer",
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
		return AddRelayerEVMCMD(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateAddRelayerFlags(cmd, args)
		if err != nil {
			return err
		}

		ProcessAddRelayerFlags(cmd, args)
		return nil
	},
}

func BindAddRelayerFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Relayer, "relayer", "", "Address to add")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	flags.MarkFlagsAsRequired(cmd, "relayer", "bridge")
}

func init() {
	BindAddRelayerFlags(addRelayerCmd)
}

func ValidateAddRelayerFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Relayer) {
		return fmt.Errorf("invalid relayer address %s", Relayer)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func ProcessAddRelayerFlags(cmd *cobra.Command, args []string) {
	RelayerAddr = common.HexToAddress(Relayer)
	BridgeAddr = common.HexToAddress(Bridge)
}

func AddRelayerEVMCMD(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	log.Debug().Msgf(`
Adding relayer
Relayer address: %s
Bridge address: %s`, Relayer, Bridge)
	_, err := contract.AddRelayer(RelayerAddr, transactor.TransactOptions{GasLimit: gasLimit})
	return err
}
