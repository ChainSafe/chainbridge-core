package bridge

import (
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
	queryProposalCmd.Flags().String("bridge", "", "bridge contract address")
	queryProposalCmd.Flags().String("dataHash", "", "hash of proposal metadata")
	queryProposalCmd.Flags().Uint64("domainId", 0, "source chain ID of proposal")
	queryProposalCmd.Flags().Uint64("depositNonce", 0, "deposit nonce of proposal")
}

func queryProposal(cmd *cobra.Command, args []string) {
	bridgeAddress := cmd.Flag("bridge").Value
	domainId := cmd.Flag("domainId").Value
	depositNonce := cmd.Flag("depositNonce").Value
	dataHash := cmd.Flag("dataHash").Value
	log.Debug().Msgf(`
Querying proposal
Chain ID: %d
Deposit nonce: %d
Data hash: %s
Bridge address: %s`, domainId, depositNonce, dataHash, bridgeAddress)
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
