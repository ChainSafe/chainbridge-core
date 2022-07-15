package bridge

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cancelProposalCmd = &cobra.Command{
	Use:   "cancel-proposal",
	Short: "Cancel an expired proposal",
	Long:  "The cancel-proposal subcommand cancels an expired proposal",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: cancelProposal,
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateCancelProposalFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindCancelProposalFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	cmd.Flags().StringVar(&DataHash, "data-hash", "", "Hash of proposal metadata")
	cmd.Flags().Uint8Var(&DomainID, "domain", 0, "Domain ID of proposal to cancel")
	cmd.Flags().Uint64Var(&DepositNonce, "deposit-nonce", 0, "Deposit nonce of proposal to cancel")
	flags.MarkFlagsAsRequired(cmd, "bridge", "data-hash", "domain", "deposit-nonce")
}

func init() {
	BindCancelProposalFlags(cancelProposalCmd)
}

func ValidateCancelProposalFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address: %s", Bridge)
	}
	return nil
}

func cancelProposal(cmd *cobra.Command, args []string) error {

	log.Debug().Msgf(`
Cancel propsal
Bridge address: %s
Chain ID: %d
Deposit nonce: %d
DataHash: %s
`, Bridge, DomainID, DepositNonce, DataHash)
	return nil
}

/*
func cancelProposal(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Uint64("gasLimit")
	gasPrice := cctx.Uint64("gasPrice")
	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	bridgeAddress, err := cliutils.DefineBridgeAddress(cctx)
	if err != nil {
		return err
	}

	domainID := cctx.Uint64("domainId")
	depositNonce := cctx.Uint64("depositNonce")
	dataHash := cctx.String("dataHash")
	dataHashBytes := utils.SliceTo32Bytes(common.Hex2Bytes(dataHash))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	err = utils.CancelProposal(ethClient, bridgeAddress, uint8(domainID), depositNonce, dataHashBytes)
	if err != nil {
		return err
	}
	log.Info().Msgf("Setting proposal with domain ID %v and deposit nonce %v status to 'Cancelled", domainID, depositNonce)
	return nil
}
*/
