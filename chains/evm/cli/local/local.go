package local

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/spf13/cobra"
)

var LocalSetupCmd = &cobra.Command{
	Use:   "local-setup",
	Short: "Local setup",
	Long:  "Locally deploy bridge and ERC20 handler contracts with preconfigured accounts and ERC20 handler",
	RunE:  localSetup,
}

// configuration
var (
	ethEndpoint1 = "http://localhost:8545"
	ethEndpoint2 = "http://localhost:8547"
	fabric1      = evmtransaction.NewTransaction
	fabric2      = evmtransaction.NewTransaction
)

func localSetup(cmd *cobra.Command, args []string) error {
	// init client1
	ethClient, err := evmclient.NewEVMClientFromParams(ethEndpoint1, EveKp.PrivateKey())
	if err != nil {
		return err
	}

	// init client2
	ethClient2, err := evmclient.NewEVMClientFromParams(ethEndpoint2, EveKp.PrivateKey())
	if err != nil {
		return err
	}

	// chain 1
	// domainsId: 0
	config, err := PrepareLocalEVME2EEnv(ethClient, fabric1, 1, big.NewInt(1), EveKp.CommonAddress())
	if err != nil {
		return err
	}

	// chain 2
	// domainId: 1
	config2, err := PrepareLocalEVME2EEnv(ethClient2, fabric2, 2, big.NewInt(1), EveKp.CommonAddress())
	if err != nil {
		return err
	}

	prettyPrint(config, config2)

	return nil
}

func prettyPrint(config, config2 EVME2EConfig) {
	fmt.Printf(`
===============================================
🎉🎉🎉 ChainBridge Successfully Deployed 🎉🎉🎉

- Chain 1 - 
Bridge: %s
ERC20: %s
ERC20 Handler: %s
Generic Handler: %s
Asset Store: %s

- Chain 2 -
Bridge: %s
ERC20: %s
ERC20 Handler: %s
Generic Handler: %s
Asset Store: %s

===============================================
`,
		// config
		config.BridgeAddr,
		config.Erc20Addr,
		config.Erc20HandlerAddr,
		config.GenericHandlerAddr,
		config.AssetStoreAddr,
		// config2
		config2.BridgeAddr,
		config2.Erc20Addr,
		config2.Erc20HandlerAddr,
		config2.GenericHandlerAddr,
		config2.AssetStoreAddr,
	)
}
