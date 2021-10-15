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
	cancelProposalCmd.Flags().StringVarP(&Bridge, "bridge", "b", "", "bridge contract address")
	cancelProposalCmd.Flags().StringVarP(&DataHash, "dataHash", "dh", "", "hash of proposal metadata")
	cancelProposalCmd.Flags().Uint64VarP(&DomainID, "domainId", "dID", 0, "domain ID of proposal to cancel")
	cancelProposalCmd.Flags().Uint64VarP(&DepositNonce, "depositNonce", "dn", 0, "deposit nonce of proposal to cancel")
}

func cancelProposal(cmd *cobra.Command, args []string) {

	log.Debug().Msgf(`
Cancel propsal
Bridge address: %s
Chain ID: %d
Deposit nonce: %d
DataHash: %s
`, Bridge, DomainID, DepositNonce, DataHash)
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
