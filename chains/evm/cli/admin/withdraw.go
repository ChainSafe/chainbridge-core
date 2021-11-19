package admin

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Withdraw tokens from a handler contract",
	Long:  "Withdraw tokens from a handler contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return WithdrawCmd(cmd, args, evmtransaction.NewTransaction, &evmgaspricer.LondonGasPriceDeterminant{})
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateWithdrawCmdFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessWithdrawCmdFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindWithdrawCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Amount, "amount", "", "token amount to withdraw. Should be set or ID or amount if both set error will occur")
	cmd.Flags().StringVar(&TokenID, "tokenId", "", "token ID to withdraw. Should be set or ID or amount if both set error will occur")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	cmd.Flags().StringVar(&Handler, "handler", "", "handler contract address")
	cmd.Flags().StringVar(&Token, "token", "", "ERC20 or ERC721 token contract address")
	cmd.Flags().StringVar(&Recipient, "recipient", "", "address to withdraw to")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	flags.MarkFlagsAsRequired(withdrawCmd, "amount", "tokenId", "bridge", "handler", "token", "recipient", "decimals")
}

func init() {
	BindWithdrawCmdFlags(withdrawCmd)
}

func ValidateWithdrawCmdFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address: %s", Bridge)
	}
	if !common.IsHexAddress(Handler) {
		return fmt.Errorf("invalid handler address: %s", Handler)
	}
	if !common.IsHexAddress(Token) {
		return fmt.Errorf("invalid token address: %s", Token)
	}
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address: %s", Recipient)
	}
	if TokenID != "" && Amount != "" {
		return errors.New("only id or amount should be set")
	}
	if TokenID == "" && Amount == "" {
		return errors.New("id or amount flag should be set")
	}
	return nil
}

func ProcessWithdrawCmdFlags(cmd *cobra.Command, args []string) error {
	var err error

	bridgeAddr = common.HexToAddress(Bridge)
	handlerAddr = common.HexToAddress(Handler)
	tokenAddr = common.HexToAddress(Token)
	recipientAddr = common.HexToAddress(Recipient)
	decimals := big.NewInt(int64(Decimals))
	realAmount, err = calls.UserAmountToWei(Amount, decimals)
	if err != nil {
		return err
	}
	return nil
}

func WithdrawCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	fmt.Printf("Withdrawing %s token from handler: %s", Amount, Handler)

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client initialization error: %v", err))
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})

	txHash, err := calls.Withdraw(
		ethClient,
		txFabric,
		gasPricer,
		gasLimit,
		bridgeAddr,
		handlerAddr,
		tokenAddr,
		recipientAddr,
		realAmount,
	)
	if err != nil {
		log.Error().Err(fmt.Errorf("admin withdrawal error: %v", err))
		return err
	}

	log.Info().Msgf("%s tokens were withdrawn from handler contract %s into recipient %s; tx hash: %s", Amount, Handler, Recipient, txHash.Hex())

	return nil
}
