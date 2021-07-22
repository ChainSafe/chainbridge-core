package deploy

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"math/big"
)

var DeployEVM = &cobra.Command{
	Use: "deploy",
	Short: "deploy smart contracts",
	Long: "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	Run: deploy,
}

func deploy(cmd *cobra.Command, args []string) {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")

	sender, err := cliutils.DefineSender(cctx)
	if err != nil {
		return err
	}
	chainID := cctx.Uint64("chainId")
	relayerThreshold := cctx.Int64("relayerThreshold")

	var relayerAddresses []common.Address
	relayerAddressesString := cctx.StringSlice("relayers")
	if len(relayerAddresses) == 0 {
		relayerAddresses = utils.DefaultRelayerAddresses
	} else {
		relayerAddresses = make([]common.Address, len(relayerAddresses))
		for i, addr := range relayerAddressesString {
			relayerAddresses[i] = common.HexToAddress(addr)
		}
	}
	var bridgeAddress common.Address
	bridgeAddressString := cctx.String("bridgeAddress")
	if common.IsHexAddress(bridgeAddressString) {
		bridgeAddress = common.HexToAddress(bridgeAddressString)
	}

	fee := cctx.String("fee")

	realFee, err := utils.UserAmountToWei(fee, big.NewInt(18))
	if err != nil {
		return err
	}

	deployments := make([]string, 0)
	if cctx.Bool("all") {
		deployments = append(deployments, []string{"bridge", "erc20Handler", "erc721Handler", "genericHandler", "erc20", "erc721"}...)
	} else {
		if cctx.Bool("bridge") {
			deployments = append(deployments, "bridge")
		}
		if cctx.Bool("erc20Handler") {
			deployments = append(deployments, "erc20Handler")
		}
		if cctx.Bool("erc721Handler") {
			deployments = append(deployments, "erc721Handler")
		}
		if cctx.Bool("genericHandler") {
			deployments = append(deployments, "genericHandler")
		}
		if cctx.Bool("erc20") {
			deployments = append(deployments, "erc20")
		}
		if cctx.Bool("erc721") {
			deployments = append(deployments, "erc721")
		}
	}
	if len(deployments) == 0 {
		return ErrNoDeploymentFalgsProvided
	}
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice), big.NewFloat(1))
	if err != nil {
		return err
	}
	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			bridgeAddress, err = utils.DeployBridge(ethClient, uint8(chainID), relayerAddresses, big.NewInt(relayerThreshold), realFee)
			if err != nil {
				return err
			}
			deployedContracts["bridge"] = bridgeAddress.String()
		case "erc20Handler":
			if bridgeAddress.String() == "" {
				return errors.New("bridge flag or bridgeAddress param should be set for contracts deployments")
			}
			erc20HandlerAddr, err := utils.DeployERC20Handler(ethClient, bridgeAddress)
			deployedContracts["erc20Handler"] = erc20HandlerAddr.String()
			if err != nil {
				return err
			}
		case "erc721Handler":
			if bridgeAddress.String() == "" {
				return errors.New("bridge flag or bridgeAddress param should be set for contracts deployments")
			}
			erc721HandlerAddr, err := utils.DeployERC721Handler(ethClient, bridgeAddress)
			deployedContracts["erc721Handler"] = erc721HandlerAddr.String()
			if err != nil {
				return err
			}
		case "genericHandler":
			if bridgeAddress.String() == "" {
				return errors.New("bridge flag or bridgeAddress param should be set for contracts deployments")
			}
			genericHandlerAddr, err := utils.DeployGenericHandler(ethClient, bridgeAddress)
			deployedContracts["genericHandler"] = genericHandlerAddr.String()
			if err != nil {
				return err
			}
		case "erc20":
			name := cctx.String("erc20Name")
			symbol := cctx.String("erc20Symbol")
			if name == "" || symbol == "" {
				return errors.New("erc20Name and erc20Symbol flags should be provided")
			}
			erc20Token, err := utils.DeployERC20Token(ethClient, name, symbol)
			deployedContracts["erc20Token"] = erc20Token.String()
			if err != nil {
				return err
			}
		case "erc721":
			erc721Token, err := utils.DeployERC721Token(ethClient)
			deployedContracts["erc721Token"] = erc721Token.String()
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("%+v", deployedContracts)
	return nil
}