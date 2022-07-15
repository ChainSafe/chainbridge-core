package admin

import (
	"errors"
	"fmt"
	"math/big"

	callsUtil "github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Withdraw tokens from a handler contract",
	Long:  "The withdraw subcommand withdrawals tokens from a handler contract",
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
		return WithdrawCmd(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateWithdrawFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessWithdrawFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindWithdrawFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Amount, "amount", "", "Token amount to withdraw, use only if ERC20 token is withdrawn. If both amount and token are set an error will occur")
	cmd.Flags().StringVar(&TokenID, "token", "", "Token ID to withdraw, use only if ERC721 token is withdrawn. If both amount and token are set an error will occur")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	cmd.Flags().StringVar(&Handler, "handler", "", "Handler contract address")
	cmd.Flags().StringVar(&Token, "token-contract", "", "ERC20 or ERC721 token contract address")
	cmd.Flags().StringVar(&Recipient, "recipient", "", "Address to withdraw to")
	cmd.Flags().Uint64Var(&Decimals, "decimals", 0, "ERC20 token decimals")
	flags.MarkFlagsAsRequired(withdrawCmd, "amount", "token", "bridge", "handler", "token-contract", "recipient", "decimals")
}

func init() {
	BindWithdrawFlags(withdrawCmd)
}

func ValidateWithdrawFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address: %s", Bridge)
	}
	if !common.IsHexAddress(Handler) {
		return fmt.Errorf("invalid handler address: %s", Handler)
	}
	if !common.IsHexAddress(Token) {
		return fmt.Errorf("invalid token-contract address: %s", Token)
	}
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address: %s", Recipient)
	}
	if TokenID != "" && Amount != "" {
		return errors.New("only token or amount should be set")
	}
	if TokenID == "" && Amount == "" {
		return errors.New("token or amount flag should be set")
	}
	return nil
}

func ProcessWithdrawFlags(cmd *cobra.Command, args []string) error {
	var err error

	BridgeAddr = common.HexToAddress(Bridge)
	HandlerAddr = common.HexToAddress(Handler)
	TokenAddr = common.HexToAddress(Token)
	RecipientAddr = common.HexToAddress(Recipient)
	decimals := big.NewInt(int64(Decimals))
	RealAmount, err = callsUtil.UserAmountToWei(Amount, decimals)
	if err != nil {
		return err
	}
	return nil
}

func WithdrawCmd(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	h, err := contract.Withdraw(
		HandlerAddr, TokenAddr, RecipientAddr, RealAmount, transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		log.Error().Err(fmt.Errorf("admin withdrawal error: %v", err))
		return err
	}

	log.Info().Msgf("%s tokens were withdrawn from handler contract %s into recipient %s; tx hash: %s", Amount, Handler, Recipient, h.Hex())
	return nil
}
