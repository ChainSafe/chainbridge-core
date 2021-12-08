package centrifuge

import (
	"errors"
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/centrifuge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	util2 "github.com/ChainSafe/chainbridge-core/chains/evm/calls/util"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/chainbridge-core/util"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var getHashCmd = &cobra.Command{
	Use:   "getHash",
	Short: "Returns if a given hash exists in asset store",
	Long:  "Checks _assetsStored map on Centrifuge asset store contract to find if asset hash exists.",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c)
		if err != nil {
			return err
		}
		return GetHashCmd(cmd, args, centrifuge.NewAssetStoreContract(c, storeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateGetHashFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessGetHashFlags(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	BindGetHashCmdFlags(getHashCmd)
}

func BindGetHashCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Hash, "hash", "", "A hash to lookup")
	cmd.Flags().StringVar(&Address, "address", "", "Centrifuge asset store contract address")
	flags.MarkFlagsAsRequired(cmd, "hash", "address")
}

func ValidateGetHashFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Address) {
		return errors.New("invalid Centrifuge asset store address")
	}

	return nil
}

func ProcessGetHashFlags(cmd *cobra.Command, args []string) error {
	storeAddr = common.HexToAddress(Address)
	byteHash = util2.SliceTo32Bytes([]byte(Hash))

	return nil
}

func GetHashCmd(cmd *cobra.Command, args []string, contract *centrifuge.AssetStoreContract) error {
	isAssetStored, err := contract.IsCentrifugeAssetStored(byteHash)
	if err != nil {
		log.Error().Err(fmt.Errorf("Checking if asset stored failed: %w", err))
		return err
	}

	log.Info().Msgf("The hash '%s' exists: %t", Hash, isAssetStored)
	return nil
}
