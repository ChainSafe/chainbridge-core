package admin

import (
	"errors"
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/util"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c)
		if err != nil {
			return err
		}
		return WithdrawCmd(cmd, args, bridge.NewBridgeContract(c, bridgeAddr, t))
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
	realAmount, err = client.UserAmountToWei(Amount, decimals)
	if err != nil {
		return err
	}
	return nil
}

func WithdrawCmd(cmd *cobra.Command, args []string, contract *bridge.BridgeContract) error {
	h, err := contract.Withdraw(handlerAddr, tokenAddr, recipientAddr, realAmount, transactor.TransactOptions{})
	if err != nil {
		log.Error().Err(fmt.Errorf("admin withdrawal error: %v", err))
		return err
	}

	log.Info().Msgf("%s tokens were withdrawn from handler contract %s into recipient %s; tx hash: %s", Amount, Handler, Recipient, h.Hex())
	return nil
}
