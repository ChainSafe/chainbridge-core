package centrifuge

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/centrifuge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a centrifuge asset store contract",
	Long:  "The deploy subcommand deploys a Centrifuge asset store contract that represents bridged Centrifuge assets",
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
		return DeployCentrifugeAssetStoreCmd(cmd, args, centrifuge.NewAssetStoreContract(c, common.Address{}, t))
	},
}

func BindDeployCmdFlags(cmd *cobra.Command) {}

func init() {
	BindDeployCmdFlags(deployCmd)
}

func DeployCentrifugeAssetStoreCmd(cmd *cobra.Command, args []string, contract *centrifuge.AssetStoreContract) error {
	assetStoreAddress, err := contract.DeployContract()
	if err != nil {
		log.Error().Err(fmt.Errorf("centrifuge asset store deploy failed: %w", err))
		return err
	}

	log.Info().Msgf("Deployed Centrifuge asset store to address: %s", assetStoreAddress.String())
	return nil
}
