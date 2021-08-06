package erc20

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint tokens on an ERC20 mintable contract",
	Long:  "Mint tokens on an ERC20 mintable contract",
	RunE:  CallMint,
}

func init() {
	mintCmd.Flags().String("amount", "", "amount to mint fee (in wei)")
	mintCmd.Flags().String("erc20Address", "", "ERC20 contract address")
	mintCmd.Flags().Uint64("decimal", 18, "ERC20 token decimals")
}

func CallMint(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return mint(cmd, args, txFabric)
}

func mint(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	amount := cmd.Flag("amount").Value.String()
	erc20Address := cmd.Flag("erc20Address").Value.String()


	decimals, err := cmd.Flags().GetUint64("decimal")
	if err != nil {
		return err
	}
	decimalsBigInt := big.NewInt(0).SetUint64(decimals)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}
	if !common.IsHexAddress(erc20Address) {
		log.Error().Err(errors.New("invalid erc20Address address"))
	}

	erc20Addr := common.HexToAddress(erc20Address)

	realAmount, err := cliutils.UserAmountToWei(amount, decimalsBigInt)
	if err != nil {
		log.Error().Err(err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	mintTokensInput, err := calls.PrepareMintTokensInput(erc20Addr, realAmount)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = calls.Transact(ethClient, txFabric, &erc20Addr, mintTokensInput, gasLimit)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Info().Msgf("%v tokens minted", amount)
	return nil
}
