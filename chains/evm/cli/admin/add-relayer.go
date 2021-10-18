package admin

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addRelayerCmd = &cobra.Command{
	Use:   "add-relayer",
	Short: "Add a new relayer",
	Long:  "Add a new relayer",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return AddRelayerEVMCMD(cmd, args, txFabric)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateAddRelayerFlags(cmd, args)
		if err != nil {
			return err
		}

		processAddRelayerFlags(cmd, args)
		return nil
	},
}

func BindAddRelayerFlags() {
	addRelayerCmd.Flags().StringVarP(&Relayer, "relayer", "r", "", "address to add")
	addRelayerCmd.Flags().StringVarP(&Bridge, "bridge", "b", "", "bridge contract address")
	flags.MarkFlagsAsRequired(addRelayerCmd, "relayer", "bridge")

}

func init() {
	BindAddRelayerFlags()
}

func validateAddRelayerFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Relayer) {
		return fmt.Errorf("invalid relayer address %s", Relayer)
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address %s", Bridge)
	}
	return nil
}

func processAddRelayerFlags(cmd *cobra.Command, args []string) {
	relayerAddr = common.HexToAddress(Relayer)
	bridgeAddr = common.HexToAddress(Bridge)
}

func AddRelayerEVMCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {

	log.Debug().Msgf(`
Adding relayer 
Relayer address: %s
Bridge address: %s`, Relayer, Bridge)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Msgf("Setting address %s as relayer on bridge %s", relayerAddr.String(), bridgeAddr.String())
	addRelayerInput, err := calls.PrepareAddRelayerInput(relayerAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &bridgeAddr, addRelayerInput, gasLimit, big.NewInt(0))
	if err != nil {
		log.Info().Msgf("%s added as relayer", relayerAddr)
		return err
	}
	return nil
}
