package deploy

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	Run:  deploy,
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
	url := cmd.Flag("url").Value.String()
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

	relayerThreshold := cmd.Flag("relayerThreshold").Value.String()

	var relayerAddresses []common.Address
	relayerAddressesString := cmd.Flag("relayers").Value.String()

	relayerAddressesSlice := strings.Split(relayerAddressesString, "")

	if len(relayerAddresses) == 0 {
		relayerAddresses = cliutils.DefaultRelayerAddresses
	} else {
		for i, addr := range relayerAddressesSlice {
			relayerAddresses[i] = common.HexToAddress(addr)
		}
	}

	var bridgeAddress common.Address
	bridgeAddressString := cmd.Flag("bridgeAddress").Value.String()
	if common.IsHexAddress(bridgeAddressString) {
		bridgeAddress = common.HexToAddress(bridgeAddressString)
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

	log.Debug().Msgf("ethClient: %v", ethClient)

	chainIdString := cmd.Flag("chainId").Value.String()

	chainIdBigInt := big.NewInt(0)

	log.Debug().Msgf(`
	Deploying
	URL: %s
	Gas limit: %d
	Gas price: %d
	Sender: %s
	Chain ID: %s
	Relayer threshold: %s
	Relayer addresses: %v`, url, gasLimit, gasPrice, sender.Address(), chainIdBigInt, relayerThreshold, relayerAddresses)

	auth, err := bind.NewKeyedTransactorWithChainID(sender.PrivateKey(), big.NewInt(0))
	if err != nil {
		log.Fatal().Err(err)
	}

	ethClient.Configure(url, auth, sender)

	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			chainIdInt, err := strconv.Atoi(chainIdString)
			if err != nil {
				log.Fatal().Err(err)
			}

			relayerThresholdInt, err := strconv.Atoi(relayerThreshold)
			if err != nil {
				log.Fatal().Err(err)
			}

			bridgeAddress, err = cliutils.DeployBridge(ethClient, auth, uint8(chainIdInt), relayerAddresses, big.NewInt(int64(relayerThresholdInt)))
			if err != nil {
				log.Fatal().Err(err)
			}
			deployedContracts["bridge"] = bridgeAddress.String()
		case "erc20Handler":
			if bridgeAddress.String() == "" {
				log.Fatal().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
			}
			erc20HandlerAddr, err := cliutils.DeployERC20Handler(ethClient, auth, bridgeAddress)
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
			erc20Token, err := cliutils.DeployERC20Token(ethClient, auth, name, symbol)
			deployedContracts["erc20Token"] = erc20Token.String()
			if err != nil {
				log.Fatal().Err(err)
			}
			if name == "" || symbol == "" {
				log.Fatal().Err(errors.New("erc20Name and erc20Symbol flags should be provided"))
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
