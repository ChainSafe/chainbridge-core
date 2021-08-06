package bridge

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

var setBurnCmd = &cobra.Command{
	Use:   "set-burn",
	Short: "Set a token contract as mintable/burnable",
	Long:  "Set a token contract as mintable/burnable in a handler",
	RunE:  func(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return SetBurnEVMCMD(cmd, args, txFabric)
},
}
func BindBridgeSetBurnCLIFlags(cli *cobra.Command) {
	cli.Flags().String("handler", "", "ERC20 handler contract address")
	cli.Flags().String("bridge", "", "bridge contract address")
	cli.Flags().String("tokenContract", "", "token contract to be registered")
}

func init() {
	BindBridgeSetBurnCLIFlags(setBurnCmd)
}

func SetBurnEVMCMD(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	handlerAddress := cmd.Flag("handler").Value.String()
	bridgeAddress := cmd.Flag("bridge").Value.String()
	tokenAddress := cmd.Flag("tokenContract").Value.String()

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	if !common.IsHexAddress(handlerAddress) {
		err := errors.New("handler address is incorrect format")
		log.Error().Err(err)
		return err
	}

	if !common.IsHexAddress(tokenAddress) {
		err := errors.New("tokenContract address is incorrect format")
		log.Error().Err(err)
		return err
	}
	handlerAddr := common.HexToAddress(handlerAddress)
	bridgeAddr := common.HexToAddress(bridgeAddress)
	tokenContractAddr := common.HexToAddress(tokenAddress)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("Setting contract %s as burnable on handler %s", tokenContractAddr.String(), handlerAddr.String())
	setBurnableInput, err := calls.PrepareSetBurnableInput(ethClient, handlerAddr, tokenContractAddr)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &bridgeAddr, setBurnableInput, gasLimit)
	if err != nil {
		log.Info().Msg("Burnable set")
		return err
	}
	return nil
}
