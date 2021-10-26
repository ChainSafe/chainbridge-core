package account

import (
	"bytes"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	accountutils "github.com/ChainSafe/chainbridge-core/keystore/account"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var importPrivKeyCmd = &cobra.Command{
	Use:   "import",
	Short: "Import bridge keystore",
	Long:  "The import subcommand is used to import a keystore for the bridge.",
	RunE:  importPrivKey,
}

func BindImportPrivKeyFlags() {
	importPrivKeyCmd.Flags().StringVarP(&PrivateKey, "privateKey", "pk", "", "Private key to encrypt")
	importPrivKeyCmd.Flags().StringVarP(&Pass, "password", "p", "", "password to encrypt with")
	flags.MarkFlagsAsRequired(importPrivKeyCmd, "privateKey", "password")
}

func init() {
	BindImportPrivKeyFlags()
}

func importPrivKey(cmd *cobra.Command, args []string) error {
	pwdb := bytes.NewBufferString(Pass)
	res, err := accountutils.ImportPrivKey(".", PrivateKey, pwdb.Bytes())
	if err != nil {
		return err
	}
	log.Debug().Msgf("filepath: %s", res)
	return nil
}
