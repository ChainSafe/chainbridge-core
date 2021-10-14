package erc721

import (
	"github.com/spf13/cobra"
)

var ERC721Cmd = &cobra.Command{
	Use:   "erc721",
	Short: "ERC721-related instructions",
	Long:  "ERC721-related instructions",
}

func init() {
	ERC721Cmd.AddCommand(mintCmd)
	ERC721Cmd.AddCommand(approveCmd)
	ERC721Cmd.AddCommand(ownerCmd)
	ERC721Cmd.AddCommand(depositCmd)
}
