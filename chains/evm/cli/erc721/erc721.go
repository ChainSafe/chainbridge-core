package erc721

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var ERC721Cmd = &cobra.Command{
	Use:   "erc721",
	Short: "ERC721-related instructions",
	Long:  "ERC721-related instructions",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// fetch global flag values
		url, gasLimit, gasPrice, senderKeyPair, err = flags.GlobalFlagValues(cmd)
		if err != nil {
			return fmt.Errorf("could not get global flags: %v", err)
		}
		return nil
	},
}

func init() {
	ERC721Cmd.AddCommand(mintCmd)
	ERC721Cmd.AddCommand(approveCmd)
	ERC721Cmd.AddCommand(ownerCmd)
	ERC721Cmd.AddCommand(depositCmd)
	ERC721Cmd.AddCommand(addMinterCmd)
}
