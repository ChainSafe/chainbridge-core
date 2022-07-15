package erc721

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var ERC721Cmd = &cobra.Command{
	Use:   "erc721",
	Short: "Set of commands for interacting with an ERC721 contract",
	Long:  "Set of commands for interacting with an ERC721 contract",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// fetch global flag values
		url, gasLimit, gasPrice, senderKeyPair, prepare, err = flags.GlobalFlagValues(cmd)
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
