package erc721

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit ERC721 token",
	Long:  "Deposit ERC721 token",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		erc721Contract, err := initializeErc721Contract()
		if err != nil {
			return err
		}
		return DepositCmd(cmd, args, erc721Contract)
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

func BindDepositCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Recipient, "recipient", "", "address of recipient")
	cmd.Flags().StringVar(&Bridge, "bridge", "", "address of bridge contract")
	cmd.Flags().StringVar(&DestionationID, "destId", "", "destination domain ID")
	cmd.Flags().StringVar(&ResourceID, "resourceId", "", "resource ID for transfer")
	cmd.Flags().StringVar(&TokenId, "tokenId", "", "ERC721 token ID")
	cmd.Flags().StringVar(&Metadata, "metadata", "", "ERC721 metadata")
	flags.MarkFlagsAsRequired(cmd, "recipient", "bridge", "destId", "resourceId", "tokenId")
}

func init() {
	BindDepositCmdFlags(depositCmd)
}

func ValidateDepositFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Recipient) {
		return fmt.Errorf("invalid recipient address")
	}
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address")
	}
	return nil
}

func ProcessDepositFlags(cmd *cobra.Command, args []string) error {
	var err error

	recipientAddr = common.HexToAddress(Recipient)
	bridgeAddr = common.HexToAddress(Bridge)

	destinationID, err = strconv.Atoi(DestionationID)
	if err != nil {
		log.Error().Err(fmt.Errorf("destination ID conversion error: %v", err))
		return err
	}

	var ok bool
	tokenId, ok = big.NewInt(0).SetString(TokenId, 10)
	if !ok {
		return fmt.Errorf("invalid token id value")
	}

	resourceId, err = flags.ProcessResourceID(ResourceID)
	return err
}

func DepositCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	// txHash, err := erc721Contract.Deposit(tokenId, Metadata, destinationID, resourceId, bridgeAddr, recipientAddr, transactor.NewDefaultTransactOptions())
	// if err != nil {
	// 	return err
	// }

	// log.Info().Msgf(
	// 	`erc721 deposit hash: %s
	// 	%s token were transferred to %s from %s`,
	// 	txHash.Hex(),
	// 	tokenId.String(),
	// 	recipientAddr.Hex(),
	// 	senderKeyPair.CommonAddress().String(),
	// )
	// return err
	return nil
}
