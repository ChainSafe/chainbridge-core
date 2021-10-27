package local

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
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
	localEndpoint1 = "http://localhost:8545"
	localEndpoint2 = "http://localhost:8546"
	localDomainId1 = 0
	localDomainId2 = 1
	fabric1        = evmtransaction.NewTransaction
	fabric2        = evmtransaction.NewTransaction
)

func localSetup(cmd *cobra.Command, args []string) error {
	// init client1
	ethClient, err := evmclient.NewEVMClientFromParams(localEndpoint1, evm.AliceKp.PrivateKey())
	if err != nil {
		return err
	}

	// init client2
	ethClient2, err := evmclient.NewEVMClientFromParams(localEndpoint2, evm.AliceKp.PrivateKey())
	if err != nil {
		return err
	}

	// chain 1
	// domainsId: 0
	bridgeAddr, erc20Addr, erc20HandlerAddr, err := evm.PrepareEVME2EEnv(ethClient, fabric1, uint8(localDomainId1), big.NewInt(1), evm.AliceKp.CommonAddress())
	if err != nil {
		return err
	}

	// chain 2
	// domainId: 1
	bridgeAddr2, erc20Addr2, erc20HandlerAddr2, err := evm.PrepareEVME2EEnv(ethClient2, fabric2, uint8(localDomainId2), big.NewInt(1), evm.AliceKp.CommonAddress())
	if err != nil {
		return err
	}

	prettyPrint(
		bridgeAddr,
		bridgeAddr2,
		erc20Addr,
		erc20Addr2,
		erc20HandlerAddr,
		erc20HandlerAddr2,
	)

	return nil
}

func prettyPrint(params ...common.Address) {
	fmt.Printf(`
===============================================
🎉🎉🎉 ChainBridge Successfully Deployed 🎉🎉🎉

- Chain 1 - 
Bridge: %s
ERC20: %s
ERC20 Handler: %s

- Chain 2 -
Bridge: %s
ERC20: %s
ERC20 Handler: %s

===============================================
`, params[0],
		params[1],
		params[2],
		params[3],
		params[4],
		params[5],
	)
}
