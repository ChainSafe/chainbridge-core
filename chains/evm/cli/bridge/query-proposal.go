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

var queryProposalCmd = &cobra.Command{
	Use:   "query-proposal",
	Short: "Query an inbound proposal",
	Long:  "The query-proposal subcommand queries an inbound proposal",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: queryProposal,
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateQueryProposalFlags(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func BindQueryProposalFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Bridge, "bridge", "", "Bridge contract address")
	cmd.Flags().StringVar(&DataHash, "data-hash", "", "Hash of proposal metadata")
	cmd.Flags().Uint8Var(&DomainID, "domain", 0, "Source domain ID of proposal")
	cmd.Flags().Uint64Var(&DepositNonce, "deposit-nonce", 0, "Deposit nonce of proposal")
	flags.MarkFlagsAsRequired(cmd, "bridge", "data-hash", "domain", "deposit-nonce")
}

func init() {
	BindQueryProposalFlags(queryProposalCmd)
}

func ValidateQueryProposalFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address: %s", Bridge)
	}
	return nil
}

func queryProposal(cmd *cobra.Command, args []string) error {
	log.Debug().Msgf(`
Querying proposal
Chain ID: %d
Deposit nonce: %d
Data hash: %s
Bridge address: %s`, DomainID, DepositNonce, DataHash, Bridge)
	return nil
}

/*
func queryProposal(cctx *cli.Context) error {
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

	prop, err := utils.QueryProposal(ethClient, bridgeAddress, uint8(domainID), depositNonce, dataHashBytes)
	if err != nil {
		return err
	}
	log.Info().Msgf("proposal with domainID %v and depositNonce %v queried. %+v", domainID, depositNonce, prop)
	return nil
}
*/
