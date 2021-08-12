package account

import (
	"bytes"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	accountutils "github.com/ChainSafe/chainbridge-core/keystore/account"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var AccountRootCMD = &cobra.Command{
	Use:   "account",
	Short: "account instructions",
	Long:  "account instructions",
}

func init() {
	AccountRootCMD.AddCommand(importPrivKeyCmd)
	AccountRootCMD.AddCommand(generateKeyPairCmd)
	BindImportPrivKeyFlags(importPrivKeyCmd)
}

var importPrivKeyCmd = &cobra.Command{
	Use:   "import",
	Short: "Import bridge keystore",
	Long:  "The import subcommand is used to import a keystore for the bridge.",
	RunE:  importPrivKey,
}

var generateKeyPairCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate bridge keystore (Secp256k1)",
	Long:  "The generate subcommand is used to generate the bridge keystore. If no options are specified, a secp256k1 key will be made.",
	RunE:  generateKeyPair,
}

func importPrivKey(cmd *cobra.Command, args []string) error {
	pk, err := cmd.Flags().GetString("privateKey")
	if err != nil {
		return err
	}
	pwd, err := cmd.Flags().GetString("password")
	if err != nil {
		return err
	}
	pwdb := bytes.NewBufferString(pwd)
	res, err := accountutils.ImportPrivKey(".", pk, pwdb.Bytes())
	if err != nil {
		return err
	}
	log.Debug().Msgf("filepath: %s", res)
	return nil
}

func generateKeyPair(cmd *cobra.Command, args []string) error {
	kp, err := secp256k1.GenerateKeypair()
	if err != nil {
		return err
	}
	log.Debug().Msgf("Addr: %s,  PrivKey %x", kp.CommonAddress().String(), kp.Encode())
	return nil
}

func BindImportPrivKeyFlags(cli *cobra.Command) {
	cli.Flags().String("privateKey", "", "Private key to encrypt")
	cli.Flags().String("password", "", "password to encrypt with")

}
