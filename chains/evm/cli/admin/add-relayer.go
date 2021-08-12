package admin

import (
	"errors"
	"fmt"
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
	RunE:   func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return AddRelayerEVMCMD(cmd, args, txFabric)
	},
}

func BindAddRelayerFlags(cli *cobra.Command) {
	cli.Flags().String("relayer", "", "address to add")
	cli.Flags().String("bridge", "", "bridge contract address")
}

func init() {
	BindAddRelayerFlags(addRelayerCmd)
}

func AddRelayerEVMCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	relayerAddress := cmd.Flag("relayer").Value.String()
	bridgeAddress := cmd.Flag("bridge").Value.String()
	log.Debug().Msgf(`
Adding relayer 
Relayer address: %s
Bridge address: %s`, relayerAddress, bridgeAddress)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(relayerAddress) {
		err := errors.New("handler address is incorrect format")
		log.Error().Err(err)
		return err
	}

	if !common.IsHexAddress(bridgeAddress) {
		err := errors.New("tokenContract address is incorrect format")
		log.Error().Err(err)
		return err
	}
	relayer := common.HexToAddress(relayerAddress)
	bridge := common.HexToAddress(bridgeAddress)
	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Msgf("Setting address %s as relayer on bridge %s", relayer.String(), bridge.String())
	addRelayerInput, err := calls.PrepareAddRelayerInput(relayer)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &bridge, addRelayerInput, gasLimit)
	if err != nil {
		log.Info().Msgf("%s added as relayer", relayerAddress)
		return err
	}
	return nil
}
