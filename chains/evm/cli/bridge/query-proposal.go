package bridge

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var QueryProposalCmd = &cobra.Command{
	Use:   "query-proposal",
	Short: "Query an inbound proposal",
	Long:  "Query an inbound proposal",
	Run:   queryProposal,
}

func init() {
	QueryProposalCmd.Flags().String("bridge", "", "bridge contract address")
	QueryProposalCmd.Flags().String("dataHash", "", "hash of proposal metadata")
	QueryProposalCmd.Flags().Uint64("chainId", 0, "source chain ID of proposal")
	QueryProposalCmd.Flags().Uint64("depositNonce", 0, "deposit nonce of proposal")
}

func queryProposal(cmd *cobra.Command, args []string) {
	bridgeAddress := cmd.Flag("bridge").Value
	chainId := cmd.Flag("chainId").Value
	depositNonce := cmd.Flag("depositNonce").Value
	dataHash := cmd.Flag("dataHash").Value
	log.Debug().Msgf(`
Querying proposal
Chain ID: %d
Deposit nonce: %d
Data hash: %s
Bridge address: %s`, chainId, depositNonce, dataHash, bridgeAddress)
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

	chainID := cctx.Uint64("chainId")
	depositNonce := cctx.Uint64("depositNonce")
	dataHash := cctx.String("dataHash")
	dataHashBytes := utils.SliceTo32Bytes(common.Hex2Bytes(dataHash))

	ethClient, err := client.NewClient(url, false, sender, big.NewInt(0).SetUint64(gasLimit), big.NewInt(0).SetUint64(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}

	prop, err := utils.QueryProposal(ethClient, bridgeAddress, uint8(chainID), depositNonce, dataHashBytes)
	if err != nil {
		return err
	}
	log.Info().Msgf("proposal with chainID %v and depositNonce %v queried. %+v", chainID, depositNonce, prop)
	return nil
}
*/
