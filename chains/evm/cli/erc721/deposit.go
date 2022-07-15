package erc721

import (
	"fmt"
	"math/big"
	"strconv"

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

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit an ERC721 token",
	Long:  "The deposit subcommand creates a new ERC721 token deposit on the bridge contract",
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
		return DepositCmd(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateDepositFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessDepositFlags(cmd, args)
		return err
	},
}

func BindDepositFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Recipient, "recipient", "", "Recipient address")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	cmd.Flags().StringVar(&DestionationID, "destination", "", "Destination domain ID")
	cmd.Flags().StringVar(&ResourceID, "resource", "", "Resource ID for transfer")
	cmd.Flags().StringVar(&Token, "token", "", "ERC721 token ID")
	cmd.Flags().StringVar(&Metadata, "metadata", "", "ERC721 token metadata")
	cmd.Flags().StringVar(&Priority, "priority", "none", "Transaction priority speed (default: medium)")
	flags.MarkFlagsAsRequired(cmd, "recipient", "bridge", "destination", "resource", "token")
}

func init() {
	BindDepositFlags(depositCmd)
}

func ValidateDepositFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address")
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address")
	}
	switch Priority {
	case "none", "slow", "medium", "fast":
		return nil
	default:
		return fmt.Errorf("invalid priority value %s, supported priorities: \"slow|medium|fast\"", Priority)
	}
}

func ProcessDepositFlags(cmd *cobra.Command, args []string) error {
	RecipientAddr = common.HexToAddress(Recipient)
	BridgeAddr = common.HexToAddress(Bridge)

	DestinationID, err = strconv.Atoi(DestionationID)
	if err != nil {
		log.Error().Err(fmt.Errorf("destination ID conversion error: %v", err))
		return err
	}

	var ok bool
	TokenId, ok = big.NewInt(0).SetString(Token, 10)
	if !ok {
		return fmt.Errorf("invalid token id value")
	}

	ResourceId, err = flags.ProcessResourceID(ResourceID)
	return err
}

func DepositCmd(cmd *cobra.Command, args []string, bridgeContract *bridge.BridgeContract) error {
	txHash, err := bridgeContract.Erc721Deposit(
		TokenId, Metadata, RecipientAddr, ResourceId, uint8(DestinationID), transactor.TransactOptions{GasLimit: gasLimit, Priority: transactor.TxPriorities[Priority]},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf(
		`erc721 deposit hash: %s
		%s token were transferred to %s from %s`,
		txHash.Hex(),
		TokenId.String(),
		RecipientAddr.Hex(),
		senderKeyPair.CommonAddress().String(),
	)
	return nil
}
