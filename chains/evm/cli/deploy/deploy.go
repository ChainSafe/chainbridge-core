package deploy

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/cliutils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var ErrNoDeploymentFalgsProvided = errors.New("provide at least one deployment flag. For help use --help.")

var DeployEVM = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	RunE:  CallDeployCLI,
}

func init() {
	flags.BindDeployEVMFlags(DeployEVM)
}

func CallDeployCLI(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return DeployCLI(cmd, args, txFabric)
}

func DeployCLI(cmd *cobra.Command, args []string, txFabric calls.TxFabric) error {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}

	gasLimit, err := cmd.Flags().GetUint64("gasLimit")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("gas limit error: %v", err))
	}
	gasPrice, err := cmd.Flags().GetUint64("gasPrice")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("gas price error: %v", err))
	}
	log.Debug().Msgf("url: %s gas limit: %v gas price: %v", url, gasLimit, gasPrice)

	senderKeyPair, err := cliutils.DefineSender(cmd)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("define sender error: %v", err))
	}

	relayersSlice, err := cmd.Flags().GetStringSlice("relayers")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("relayers error: %v", err))
	}

	relayerThreshold, err := cmd.Flags().GetUint64("relayerThreshold")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("relayer threshold error: %v", err))
	}
	log.Debug().Msg("got relayer threshold")

	relayerAddressesStringSlice := viper.GetStringSlice(flags.RelayersFlagName)
	log.Debug().Msgf("relayer addresses from viper: %v", relayerAddressesStringSlice)

	log.Debug().Msgf("relayers: %s", relayersSlice)
	log.Debug().Msgf("relayers count: %d", len(relayersSlice))

	var relayerAddresses []common.Address
	for _, addr := range relayersSlice {
		relayerAddresses = append(relayerAddresses, common.HexToAddress(addr))
	}
	log.Debug().Msg("got relayers")

	var bridgeAddress common.Address
	bridgeAddressString := cmd.Flag("bridgeAddress").Value.String()
	if common.IsHexAddress(bridgeAddressString) {
		bridgeAddress = common.HexToAddress(bridgeAddressString)
	}
	log.Debug().Msg("got bridge address")

	deployments := make([]string, 0)

	// flag bools
	log.Debug().Msgf("all bool: %v", viper.GetBool("all"))

	allBool, err := cmd.Flags().GetBool("all")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("all flag error: %v", err))
	}

	bridgeBool, err := cmd.Flags().GetBool("bridge")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("bridge flag error: %v", err))
	}

	erc20HandlerBool, err := cmd.Flags().GetBool("erc20Handler")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("erc20 handler flag error: %v", err))
	}

	erc721HandlerBool, err := cmd.Flags().GetBool("erc721Handler")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("erc721 handler flag error: %v", err))
	}

	genericHandlerBool, err := cmd.Flags().GetBool("genericHandler")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("generic handler flag error: %v", err))
	}

	erc20Bool, err := cmd.Flags().GetBool("erc20")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("erc20 flag error: %v", err))
	}

	erc721Bool, err := cmd.Flags().GetBool("erc721")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("erc721 flag error: %v", err))
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

	chainId := cmd.Flag("chainId").Value.String()

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Fatal().Err(fmt.Errorf("eth client intialization error: %v", err))
	}

	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			log.Debug().Msgf("deploying bridge..")
			// convert chain ID to uint
			chainIdInt, err := strconv.Atoi(chainId)
			if err != nil {
				log.Fatal().Err(fmt.Errorf("chain ID flag error: %v", err))
			}
			bridgeAddr, err := calls.DeployBridge(ethClient, txFabric, uint8(chainIdInt), relayerAddresses, big.NewInt(0).SetUint64(relayerThreshold))
			if err != nil {
				log.Fatal().Err(fmt.Errorf("bridge deploy failed: %w", err))
			}
			deployedContracts["bridge"] = bridgeAddr.String()

			log.Debug().Msgf("bridge address; %v", bridgeAddr.String())
		case "erc20Handler":
			log.Debug().Msgf("deploying ERC20 handler..")
			if bridgeAddress.String() == "" {
				log.Fatal().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
			}
			erc20HandlerAddr, err := calls.DeployErc20Handler(ethClient, txFabric, bridgeAddress)
			if err != nil {
				log.Fatal().Err(fmt.Errorf("ERC20 handler deploy failed: %w", err))
			}
			deployedContracts["erc20Handler"] = erc20HandlerAddr.String()
		case "erc20":
			log.Debug().Msgf("deploying ERC20..")
			name := cmd.Flag("erc20Name").Value.String()
			symbol := cmd.Flag("erc20Symbol").Value.String()
			if name == "" || symbol == "" {
				log.Fatal().Err(errors.New("erc20Name and erc20Symbol flags should be provided"))
			}

			erc20Addr, err := calls.DeployErc20(ethClient, txFabric, name, symbol)
			if err != nil {
				log.Fatal().Err(fmt.Errorf("erc 20 deploy failed: %w", err))
			}
			deployedContracts["erc20Token"] = erc20Addr.String()
			if err != nil {
				log.Fatal().Err(err)
			}
			if name == "" || symbol == "" {
				log.Fatal().Err(errors.New("erc20Name and erc20Symbol flags should be provided"))
			}
			//case "erc721":
			//	log.Debug().Msgf("deploying ERC721..")
			//	erc721Token, err := cliutils.DeployERC721Token(ethClient)
			//	deployedContracts["erc721Token"] = erc721Token.String()
			//	if err != nil {
			//		log.Fatal().Err(err)
			//	}
			//case "erc721Handler":
			//	log.Debug().Msgf("deploying ERC721 handler..")
			//	if bridgeAddress.String() == "" {
			//		log.Fatal().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
			//	}
			//	erc721HandlerAddr, err := cliutils.DeployERC721Handler(ethClient, bridgeAddress)
			//	deployedContracts["erc721Handler"] = erc721HandlerAddr.String()
			//	if err != nil {
			//		log.Fatal().Err(err)
			//	}
			//case "genericHandler":
			//	log.Debug().Msgf("deploying generic handler..")
			//	if bridgeAddress.String() == "" {
			//		log.Fatal().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
			//	}
			//	genericHandlerAddr, err := cliutils.DeployGenericHandler(ethClient, bridgeAddress)
			//	deployedContracts["genericHandler"] = genericHandlerAddr.String()
			//	if err != nil {
			//		log.Fatal().Err(err)
			//	}
		}
	}
	fmt.Printf("%+v", deployedContracts)
	return nil
}
