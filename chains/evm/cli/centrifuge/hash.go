package centrifuge

import (
	"errors"
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var hashExistsCmd = &cobra.Command{
	Use:   "hash-exists",
	Short: "Return if a given hash exists in asset store",
	Long:  "Calls ",
	RunE: func(cmd *cobra.Command, args []string) error {
		txFabric := evmtransaction.NewTransaction
		return HashExistsCmd(cmd, args, txFabric)
	},
}

func BindHashExistsCmdFlags(cli *cobra.Command) {
	cli.Flags().String("hash", "", "A hash to lookup")
	cli.Flags().String("address", "", "Centrifuge asset store contract address")

	err := cli.MarkFlagRequired("hash")
	if err != nil {
		panic(err)
	}
	err = cli.MarkFlagRequired("address")
	if err != nil {
		panic(err)
	}
}

func init() {
	BindHashExistsCmdFlags(hashExistsCmd)
}

func HashExistsCmd(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	storeAddrStr := cmd.Flag("address").Value.String()
	hash := cmd.Flag("hash").Value.String()

	if !common.IsHexAddress(storeAddrStr) {
		return errors.New("invalid Centrifuge asset store address")
	}
	storeAddr := common.HexToAddress(storeAddrStr)
	byteHash := calls.SliceTo32Bytes([]byte(hash))

	url, _, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return fmt.Errorf("could not get global flags: %v", err)
	}

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey(), gasPrice)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return err
	}

	isAssetStored, err := calls.IsAssetStored(ethClient, storeAddr, byteHash)
	if err != nil {
		log.Error().Err(fmt.Errorf("Centrifuge asset store deploy failed: %w", err))
	}

	log.Info().Msgf("The hash '%s' exists: %t", hash, isAssetStored)
	return nil
}
