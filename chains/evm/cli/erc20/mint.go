package erc20

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint tokens on an ERC20 mintable contract",
	Long:  "Mint tokens on an ERC20 mintable contract",
	Run:   mint,
}

func init() {
	mintCmd.Flags().String("amount", "", "amount to mint fee (in wei)")
	mintCmd.Flags().String("erc20Address", "", "ERC20 contract address")
}

func mint(cmd *cobra.Command, args []string) {
	amount := cmd.Flag("amount").Value.String()
	erc20Address := cmd.Flag("erc20Address").Value.String()
	log.Debug().Msgf(`
Minting token
Amount: %s
ERC20 address: %s`, amount, erc20Address)

	url := cmd.Flag("url").Value.String()
	decimals := "2"
	decimalsBigInt, _ := big.NewInt(0).SetString(decimals, 10)

	if !common.IsHexAddress(erc20Address) {
		log.Fatal().Err(errors.New("invalid erc20Address address"))
	}

	erc20Addr := common.HexToAddress(erc20Address)

	realAmount, err := cliutils.UserAmountToWei(amount, decimalsBigInt)
	if err != nil {
		log.Fatal().Err(err)
	}

	senderKeyPair, err := cliutils.DefineSender(cmd)
	if err != nil {
		log.Fatal().Err(err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Fatal().Err(err)
	}

	mintTokensInput, err := calls.PrepareMintTokensInput(erc20Addr, realAmount)
	if err != nil {
		log.Fatal().Err(err)
	}

	_, err = calls.SendInput(ethClient, erc20Addr, mintTokensInput)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("%v tokens minted", amount)
}
