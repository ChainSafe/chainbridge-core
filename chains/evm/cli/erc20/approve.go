package erc20

import (
	"errors"
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

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve tokens in an ERC20 contract for transfer",
	Long:  "Approve tokens in an ERC20 contract for transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return ApproveCmd(cmd, args, txFabric)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateApproveFlags(cmd, args)
		if err != nil {
			return err
		}

		err = processApproveFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindApproveCmdFlags() {
	balanceCmd.Flags().StringVarP(&Erc20Address, "erc20Address", "erc20add", "", "ERC20 contract address")
	depositCmd.Flags().StringVarP(&Amount, "amount", "a", "", "amount to grant allowance")
	depositCmd.Flags().StringVarP(&Recipient, "recipient", "r", "", "address of recipient")
	depositCmd.Flags().Uint64VarP(&Decimals, "decimals", "r", 18, "ERC20 token decimals")
	flags.CheckRequiredFlags(depositCmd, "erc20Address", "amount", "recipient")
}

func init() {
	BindApproveCmdFlags()
}

func validateApproveFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc20Address) {
		return errors.New("invalid erc20Address address")
	}
	if !common.IsHexAddress(Recipient) {
		return errors.New("invalid minter address")
	}
	return nil
}

func processApproveFlags(cmd *cobra.Command, args []string) error {
	var err error

	decimals := big.NewInt(int64(Decimals))
	erc20Addr = common.HexToAddress(Erc20Address)
	recipientAddress = common.HexToAddress(Recipient)
	realAmount, err = calls.UserAmountToWei(Amount, decimals)
	if err != nil {
		return err
	}

	return nil
}

func ApproveCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	log.Debug().Msgf(`
Approving ERC20
ERC20 address: %s
Recipient address: %s
Amount: %s
Decimals: %v`,
		Erc20Address, Recipient, Amount, Decimals)

	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}
	i, err := calls.PrepareErc20ApproveInput(recipientAddress, realAmount)
	if err != nil {
		log.Fatal().Err(err)
		return err
	}
	_, err = calls.Transact(ethClient, txFabric, &erc20Addr, i, gasLimit)
	if err != nil {
		log.Fatal().Err(err)
		return err
	}
	log.Info().Msgf("%s account granted allowance on %v tokens of %s", recipientAddress.String(), Amount, recipientAddress.String())
	return nil
}
