package bridge

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var setBurnCmd = &cobra.Command{
	Use:   "set-burn",
	Short: "Set a token contract as mintable/burnable",
	Long:  "Set a token contract as mintable/burnable in a handler",
	Run:   setBurn,
}

func init() {
	setBurnCmd.Flags().String("handler", "", "ERC20 handler contract address")
	setBurnCmd.Flags().String("bridge", "", "bridge contract address")
	setBurnCmd.Flags().String("tokenContract", "", "token contract to be registered")
}

func setBurn(cmd *cobra.Command, args []string) {
	handlerAddress := cmd.Flag("handler").Value
	bridgeAddress := cmd.Flag("bridge").Value
	tokenAddress := cmd.Flag("tokenContract").Value
	gasPrice, err := cmd.Flags().GetUint64("gasPrice")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("gas price error: %v", err))
	}
	log.Debug().Msgf(`
Setting contract as mintable/burnable
Handler address: %s
Bridge address: %s
Token contract address: %s`, handlerAddress, bridgeAddress, tokenAddress)

	url := cmd.Flag("url").Value.String()
	handler := cmd.Flag("handler").Value.String()
	if !common.IsHexAddress(handler) {
		log.Fatal().Err(errors.New("handler address is incorrect format"))
	}
	tokenContract := cmd.Flag("tokenContract").Value.String()
	if !common.IsHexAddress(tokenContract) {
		log.Fatal().Err(errors.New("tokenContract address is incorrect format"))
	}
	handlerAddr := common.HexToAddress(handler)
	bridgeAddr := common.HexToAddress(bridgeAddress.String())
	tokenContractAddr := common.HexToAddress(tokenContract)

	senderKeyPair, err := cliutils.DefineSender(cmd)
	if err != nil {
		log.Fatal().Err(err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), big.NewInt(0).SetUint64(gasPrice))
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Setting contract %s as burnable on handler %s", tokenContractAddr.String(), handlerAddress.String())
	setBurnableInput, err := calls.PrepareSetBurnableInput(ethClient, bridgeAddr, handlerAddr, tokenContractAddr)
	if err != nil {
		log.Fatal().Err(err)
	}

	_, err = calls.SendInput(ethClient, handlerAddr, setBurnableInput)
	if err != nil {
		log.Info().Msg("Burnable set")
	}
}
