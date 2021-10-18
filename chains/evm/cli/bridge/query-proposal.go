package bridge

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var queryProposalCmd = &cobra.Command{
	Use:   "query-proposal",
	Short: "Query an inbound proposal",
	Long:  "Query an inbound proposal",
	Run:   queryProposal,
}

func init() {
	queryProposalCmd.Flags().StringVarP(&Bridge, "bridge", "b", "", "bridge contract address")
	queryProposalCmd.Flags().StringVarP(&DataHash, "dataHash", "dh", "", "hash of proposal metadata")
	queryProposalCmd.Flags().Uint64VarP(&DomainID, "domainId", "dID", 0, "source domain ID of proposal")
	queryProposalCmd.Flags().Uint64VarP(&DepositNonce, "depositNonce", "dn", 0, "	deposit nonce of proposal")
	flags.CheckRequiredFlags(queryProposalCmd, "bridge", "dataHash", "domainId", "depositNonce")
}

func queryProposal(cmd *cobra.Command, args []string) {
	log.Debug().Msgf(`
Querying proposal
Chain ID: %d
Deposit nonce: %d
Data hash: %s
Bridge address: %s`, DomainID, DepositNonce, DataHash, Bridge)
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
