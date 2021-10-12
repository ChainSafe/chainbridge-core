package account

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	accountutils "github.com/ChainSafe/chainbridge-core/keystore/account"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var AccountRootCMD = &cobra.Command{
	Use:   "accounts",
	Short: "Account instructions",
	Long:  "Account instructions",
}

func init() {
	AccountRootCMD.AddCommand(importPrivKeyCmd)
	AccountRootCMD.AddCommand(generateKeyPairCmd)
	AccountRootCMD.AddCommand(transferBaseCurrencyCmd)
	BindImportPrivKeyFlags(importPrivKeyCmd)
	BindTransferCmdFlags(transferBaseCurrencyCmd)
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
	Long:  "The generate subcommand is used to generate the bridge keystore. If no options are specified, a Secp256k1 key will be made.",
	RunE:  generateKeyPair,
}

var transferBaseCurrencyCmd = &cobra.Command{
	Use:    "transfer",
	Short:  "Transfer base currency",
	Long:   "The generate subcommand is used to transfer the base currency",
	PreRun: confirmTransfer,
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return transferBaseCurrency(cmd, args, txFabric)
	},
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

func transferBaseCurrency(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	recipient := cmd.Flag("recipient").Value.String()
	amount := cmd.Flag("amount").Value.String()
	if !common.IsHexAddress(recipient) {
		return fmt.Errorf("invalid recipient address %s", recipient)
	}
	recipientAddress := common.HexToAddress(recipient)
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	decimals, _ := big.NewInt(0).SetString(cmd.Flag("decimals").Value.String(), 10)

	weiAmount, err := calls.UserAmountToWei(amount, decimals)
	if err != nil {
		return err
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	txHash, err := calls.Transact(ethClient, txFabric, &recipientAddress, nil, gasLimit, weiAmount)
	if err != nil {
		log.Error().Err(fmt.Errorf("base currency deposit error: %v", err))
		return err
	}

	log.Debug().Msgf("base currency transaction hash: %s", txHash.Hex())

	log.Info().Msgf("%s tokens were transferred to %s from %s", amount, recipientAddress.Hex(), senderKeyPair.CommonAddress().String())
	return nil
}

func confirmTransfer(cmd *cobra.Command, args []string) {
	recipient := cmd.Flag("recipient").Value.String()
	amount := cmd.Flag("amount").Value.String()
	decimals := cmd.Flag("decimals").Value.String()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Send transaction %s(%s) to %s (Y/N)?", amount, decimals, recipient)
		s, _ := reader.ReadString('\n')

		s = strings.ToLower(strings.TrimSuffix(s, "\n"))

		if strings.Compare(s, "n") == 0 {
			os.Exit(0)
		} else if strings.Compare(s, "y") == 0 {
			break
		} else {
			continue
		}
	}
}

func BindImportPrivKeyFlags(cli *cobra.Command) {
	cli.Flags().String("privateKey", "", "Private key to encrypt")
	cli.Flags().String("password", "", "password to encrypt with")
}

func BindTransferCmdFlags(cli *cobra.Command) {
	cli.Flags().String("recipient", "", "recipient address")
	cli.Flags().String("amount", "", "transfer amount")
	cli.Flags().Uint64("decimals", 18, "base token decimals")
	err := cli.MarkFlagRequired("amount")
	if err != nil {
		panic(err)
	}
	err = cli.MarkFlagRequired("recipient")
	if err != nil {
		panic(err)
	}
}
