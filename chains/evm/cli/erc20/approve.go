package erc20

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve tokens in an ERC20 contract for transfer",
	Long:  "Approve tokens in an ERC20 contract for transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return ApproveCmd(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
	},
}

func BindApproveCmdFlags(cli *cobra.Command) {
	cli.Flags().String("erc20address", "", "ERC20 contract address")
	cli.Flags().String("amount", "", "amount to grant allowance")
	cli.Flags().String("recipient", "", "address of recipient")
	cli.Flags().Uint64("decimals", 18, "ERC20 token decimals")
	err := cli.MarkFlagRequired("decimals")
	if err != nil {
		panic(err)
	}
}

func init() {
	BindApproveCmdFlags(approveCmd)
}

func ApproveCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	erc20Address := cmd.Flag("erc20address").Value.String()
	recipientAddress := cmd.Flag("recipient").Value.String()
	amount := cmd.Flag("amount").Value.String()
	decimals, err := cmd.Flags().GetUint64("decimals")
	if err != nil {
		return err
	}
	log.Debug().Msgf(`
Approving ERC20
ERC20 address: %s
Recipient address: %s
Amount: %s
Decimals: %v`,
		erc20Address, recipientAddress, amount, decimals)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	decimalsBigInt := big.NewInt(0).SetUint64(decimals)

	if !common.IsHexAddress(erc20Address) {
		return errors.New("invalid erc20Address address")
	}
	erc20Addr := common.HexToAddress(erc20Address)

	if !common.IsHexAddress(recipientAddress) {
		return errors.New("invalid minter address")
	}
	recipientAddr := common.HexToAddress(recipientAddress)

	realAmount, err := utils.UserAmountToWei(amount, decimalsBigInt)
	if err != nil {
		log.Fatal().Err(err)
		return err
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})
	i, err := calls.PrepareErc20ApproveInput(recipientAddr, realAmount)
	if err != nil {
		log.Fatal().Err(err)
		return err
	}
	_, err = calls.Transact(ethClient, txFabric, gasPricer, &erc20Addr, i, gasLimit, big.NewInt(0))
	if err != nil {
		log.Fatal().Err(err)
		return err
	}
	log.Info().Msgf("%s account granted allowance on %v tokens of %s", recipientAddr.String(), amount, recipientAddr.String())
	return nil
}
