package bridge

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cancelProposalCmd = &cobra.Command{
	Use:   "cancel-proposal",
	Short: "Cancel an expired proposal",
	Long:  "Cancel an expired proposal",
	Run:   cancelProposal,
}

func init() {
	cancelProposalCmd.Flags().String("bridge", "", "bridge contract address")
	cancelProposalCmd.Flags().String("dataHash", "", "hash of proposal metadata")
	cancelProposalCmd.Flags().Uint64("domainId", 0, "chain ID of proposal to cancel")
	cancelProposalCmd.Flags().Uint64("depositNonce", 0, "deposit nonce of proposal to cancel")
}

func cancelProposal(cmd *cobra.Command, args []string) {
	adminAddress := cmd.Flag("admin").Value
	bridgeAddress := cmd.Flag("bridge").Value
	domainId := cmd.Flag("domainId").Value
	depositNonce := cmd.Flag("depositNonce").Value
	log.Debug().Msgf(`
Cancel propsal
Admin address: %s
Bridge address: %s
Chain ID: %d
Deposit nonce: %d`, adminAddress, bridgeAddress, domainId, depositNonce)
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
	log.Info().Msgf("Setting proposal with chain ID %v and deposit nonce %v status to 'Cancelled", domainID, depositNonce)
	return nil
}
*/
