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
	Use:   "import-pk",
	Short: "Add a minter to an Erc20 mintable contract",
	Long:  "Add a minter to an Erc20 mintable contract",
	RunE:  importPrivKey,
}

var generateKeyPairCmd = &cobra.Command{
	Use:   "generate-keypair",
	Short: "Add a minter to an Erc20 mintable contract",
	Long:  "Add a minter to an Erc20 mintable contract",
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