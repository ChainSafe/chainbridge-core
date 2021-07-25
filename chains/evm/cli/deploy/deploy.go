package deploy

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var ErrNoDeploymentFalgsProvided = errors.New("provide at least one deployment flag. For help use --help.")

var DeployEVM = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	Run:   deploy,
}

func init() {
	DeployEVM.Flags().Bool("bridge", false, "deploy bridge")
	DeployEVM.Flags().Bool("erc20Handler", false, "deploy ERC20 handler")
	DeployEVM.Flags().Bool("erc721Handler", false, "deploy ERC721 handler")
	DeployEVM.Flags().Bool("genericHandler", false, "deploy generic handler")
	DeployEVM.Flags().Bool("erc20", false, "deploy ERC20")
	DeployEVM.Flags().Bool("erc721", false, "deploy ERC721")
	DeployEVM.Flags().Bool("all", false, "deploy all")
	DeployEVM.Flags().Int64("relayerThreshold", 1, "number of votes required for a proposal to pass")
	DeployEVM.Flags().String("chainId", "1", "chain ID for the instance")
	DeployEVM.Flags().StringSlice("relayers", []string{}, "list of initial relayers")
	DeployEVM.Flags().String("fee", "0", "fee to be taken when making a deposit (in ETH, decimas are allowed)")
	DeployEVM.Flags().String("bridgeAddress", "", "bridge contract address. Should be provided if handlers are deployed separately")
	DeployEVM.Flags().String("erc20Symbol", "", "ERC20 contract symbol")
	DeployEVM.Flags().String("erc20Name", "", "ERC20 contract name")
}

func deploy(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value
	// gasLimit := cmd.Flag("gasLimit").Value
	gasLimit, err := cmd.Flags().GetUint64("gasLimit")
	if err != nil {
		log.Fatal().Err(err)
	}

	gasPrice, err := cmd.Flags().GetUint64("gasPrice")
	if err != nil {
		log.Fatal().Err(err)
	}

	privateKey := cliutils.AliceKp.PrivateKey()

	privateKeyString := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

	sender, err := cliutils.DefineSender(privateKeyString)
	if err != nil {
		log.Fatal().Err(err)
	}

	chainId := cmd.Flag("chainId").Value
	relayerThreshold := cmd.Flag("relayerThreshold").Value

	var relayerAddresses []common.Address
	relayerAddressesString := cmd.Flag("relayers").Value

	relayerAddressesSlice := strings.Split(relayerAddressesString.String(), "")

	if len(relayerAddresses) == 0 {
		relayerAddresses = cliutils.DefaultRelayerAddresses
	} else {
		for i, addr := range relayerAddressesSlice {
			relayerAddresses[i] = common.HexToAddress(addr)
		}
	}

	log.Debug().Msgf(`
	Deploying
	URL: %s
	Gas limit: %d
	Gas price: %d
	Sender: %s
	Chain ID: %s
	Relayer threshold: %s
	Relayer addresses: %v`, url, gasLimit, gasPrice, sender.Address(), chainId, relayerThreshold, relayerAddresses)

	var bridgeAddress common.Address
	bridgeAddressString := cmd.Flag("bridgeAddress").Value.String()
	if common.IsHexAddress(bridgeAddressString) {
		bridgeAddress = common.HexToAddress(bridgeAddressString)
	}

	fee := cmd.Flag("fee").Value

	realFee, err := cliutils.UserAmountToWei(fee.String(), big.NewInt(18))
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	deployments := make([]string, 0)

	// flag bools

	allBool, err := cmd.Flags().GetBool("all")
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	bridgeBool, err := cmd.Flags().GetBool("bridge")
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	erc20HandlerBool, err := cmd.Flags().GetBool("erc20Handler")
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	erc721HandlerBool, err := cmd.Flags().GetBool("erc721Handler")
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	genericHandlerBool, err := cmd.Flags().GetBool("genericHandler")
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	erc20Bool, err := cmd.Flags().GetBool("erc20")
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	erc721Bool, err := cmd.Flags().GetBool("erc721")
	if err != nil {
		// log.Err(err)
		log.Fatal().Err(err)
	}

	if allBool {
		deployments = append(deployments, []string{"bridge", "erc20Handler", "erc721Handler", "genericHandler", "erc20", "erc721"}...)
	} else {
		if bridgeBool {
			deployments = append(deployments, "bridge")
		}
		if erc20HandlerBool {
			deployments = append(deployments, "erc20Handler")
		}
		if erc721HandlerBool {
			deployments = append(deployments, "erc721Handler")
		}
		if genericHandlerBool {
			deployments = append(deployments, "genericHandler")
		}
		if erc20Bool {
			deployments = append(deployments, "erc20")
		}
		if erc721Bool {
			deployments = append(deployments, "erc721")
		}
	}
	if len(deployments) == 0 {
		log.Fatal().Err(ErrNoDeploymentFalgsProvided)
	}
	ethClient := evmclient.NewEVMClient()
	// err = ethClient.ConfigureNoFile(url.String(), false, sender, big.NewInt(int64(gasLimit)), big.NewInt(int64(gasPrice)), big.NewFloat(1))
	// if err != nil {
	// 	log.Fatal().Err(err)
	// }

	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			chainIdInt, err := strconv.Atoi(chainId.String())
			if err != nil {
				log.Fatal().Err(err)
			}

			relayerThresholdInt, err := strconv.Atoi(relayerThreshold.String())
			if err != nil {
				log.Fatal().Err(err)
			}

			bridgeAddress, err = cliutils.DeployBridge(ethClient, uint8(chainIdInt), relayerAddresses, big.NewInt(int64(relayerThresholdInt)), realFee)
			if err != nil {
				log.Fatal().Err(err)
			}
			deployedContracts["bridge"] = bridgeAddress.String()
		case "erc20Handler":
			if bridgeAddress.String() == "" {
				log.Fatal().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
			}
			erc20HandlerAddr, err := cliutils.DeployERC20Handler(ethClient, bridgeAddress)
			deployedContracts["erc20Handler"] = erc20HandlerAddr.String()
			if err != nil {
				log.Fatal().Err(err)
			}
		case "erc721Handler":
			if bridgeAddress.String() == "" {
				log.Fatal().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
			}
			erc721HandlerAddr, err := cliutils.DeployERC721Handler(ethClient, bridgeAddress)
			deployedContracts["erc721Handler"] = erc721HandlerAddr.String()
			if err != nil {
				log.Fatal().Err(err)
			}
		case "genericHandler":
			if bridgeAddress.String() == "" {
				log.Fatal().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
			}
			genericHandlerAddr, err := cliutils.DeployGenericHandler(ethClient, bridgeAddress)
			deployedContracts["genericHandler"] = genericHandlerAddr.String()
			if err != nil {
				log.Fatal().Err(err)
			}
		case "erc20":
			name := cmd.Flag("erc20Name").Value.String()
			symbol := cmd.Flag("erc20Symbol").Value.String()
			if name == "" || symbol == "" {
				log.Fatal().Err(errors.New("erc20Name and erc20Symbol flags should be provided"))
			}
			erc20Token, err := cliutils.DeployERC20Token(ethClient, name, symbol)
			deployedContracts["erc20Token"] = erc20Token.String()
			if err != nil {
				log.Fatal().Err(err)
			}
		case "erc721":
			erc721Token, err := cliutils.DeployERC721Token(ethClient)
			deployedContracts["erc721Token"] = erc721Token.String()
			if err != nil {
				log.Fatal().Err(err)
			}
		}
	}
	fmt.Printf("%+v", deployedContracts)
}

/*
func deploy(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value
	gasLimit := cmd.Flag("gasLimit").Value
	gasPrice := cmd.Flag("gasPrice").Value
	log.Debug().Msgf("%v %v %v", url, gasLimit, gasPrice)

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
*/
