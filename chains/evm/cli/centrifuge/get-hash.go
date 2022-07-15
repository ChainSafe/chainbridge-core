package centrifuge

import (
	"errors"
	"fmt"

	callsUtil "github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/centrifuge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/util"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var getHashCmd = &cobra.Command{
	Use:   "get-hash",
	Short: "Returns the status of whether a given hash exists in an asset store",
	Long:  "The get-hash subcommand checks the _assetsStored map on a Centrifuge asset store contract to determine whether the asset hash exists or not",
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
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c, prepare)
		if err != nil {
			return err
		}
		return GetHashCmd(cmd, args, centrifuge.NewAssetStoreContract(c, StoreAddr, t))
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
	BindGetHashFlags(getHashCmd)
}

func BindGetHashFlags(cmd *cobra.Command) {
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
	StoreAddr = common.HexToAddress(Address)
	ByteHash = callsUtil.SliceTo32Bytes([]byte(Hash))

	return nil
}

func GetHashCmd(cmd *cobra.Command, args []string, contract *centrifuge.AssetStoreContract) error {
	isAssetStored, err := contract.IsCentrifugeAssetStored(ByteHash)
	if err != nil {
		log.Error().Err(fmt.Errorf("checking if asset stored failed: %w", err))
		return err
	}

	log.Info().Msgf("The hash '%s' exists: %t", Hash, isAssetStored)
	return nil
}
