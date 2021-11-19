package deploy

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var ErrNoDeploymentFlagsProvided = errors.New("provide at least one deployment flag. For help use --help")

var DeployEVM = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy smart contracts",
	Long:  "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	RunE: CallDeployCLI,
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateDeployFlags(cmd, args)
		if err != nil {
			return err
		}
		err = ProcessDeployFlags(cmd, args)
		return err
	},
}

var (
	// Flags for all EVM Deploy CLI commands
	Bridge           bool
	BridgeAddress    string
	DeployAll        bool
	DomainId         uint8
	GenericHandler   bool
	Erc20            bool
	Erc20Handler     bool
	Erc20Name        string
	Erc20Symbol      string
	Erc721           bool
	Erc721Handler    bool
	Erc721Name       string
	Erc721Symbol     string
	Erc721BaseURI    string
	Fee              uint64
	RelayerThreshold uint64
	Relayers         []string
)

func BindDeployEVMFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&Bridge, "bridge", false, "deploy bridge")
	cmd.Flags().StringVar(&BridgeAddress, "bridgeAddress", "", "bridge contract address. Should be provided if handlers are deployed separately")
	cmd.Flags().BoolVar(&DeployAll, "all", false, "deploy all")
	cmd.Flags().Uint8Var(&DomainId, "domainId", 1, "domain ID for the instance")
	cmd.Flags().BoolVar(&Erc20, "erc20", false, "deploy ERC20")
	cmd.Flags().BoolVar(&Erc20Handler, "erc20Handler", false, "deploy ERC20 handler")
	cmd.Flags().StringVar(&Erc20Name, "erc20Name", "", "ERC20 contract name")
	cmd.Flags().StringVar(&Erc20Symbol, "erc20Symbol", "", "ERC20 contract symbol")
	cmd.Flags().BoolVar(&Erc721, "erc721", false, "deploy ERC721")
	cmd.Flags().BoolVar(&Erc721Handler, "erc721Handler", false, "deploy ERC721 handler")
	cmd.Flags().StringVar(&Erc721Name, "erc721Name", "", "ERC721 contract name")
	cmd.Flags().StringVar(&Erc721Symbol, "erc721Symbol", "", "ERC721 contract symbol")
	cmd.Flags().StringVar(&Erc721BaseURI, "erc721BaseURI", "", "ERC721 base URI")
	cmd.Flags().BoolVar(&GenericHandler, "genericHandler", false, "deploy generic handler")
	cmd.Flags().Uint64Var(&Fee, "fee", 0, "fee to be taken when making a deposit (in ETH, decimas are allowed)")
	cmd.Flags().StringSliceVar(&Relayers, "relayers", []string{}, "list of initial relayers")
	cmd.Flags().Uint64Var(&RelayerThreshold, "relayerTreshold", 1, "number of votes required for a proposal to pass")
}

func init() {
	BindDeployEVMFlags(DeployEVM)
}

func ValidateDeployFlags(cmd *cobra.Command, args []string) error {
	deployments = make([]string, 0)
	if DeployAll {
		flags.MarkFlagsAsRequired(cmd, "relayerThreshold", "domainId", "fee", "erc20Symbol", "erc20Name")
		deployments = append(deployments, []string{"bridge", "erc20Handler", "erc721Handler", "genericHandler", "erc20", "erc721"}...)
	} else {
		if Bridge {
			flags.MarkFlagsAsRequired(cmd, "relayerThreshold", "domainId", "fee")
			deployments = append(deployments, "bridge")
		}
		if Erc20Handler {
			if !Bridge {
				flags.MarkFlagsAsRequired(cmd, "bridgeAddress")
			}
			deployments = append(deployments, "erc20Handler")
		}
		if Erc721Handler {
			if !Bridge {
				flags.MarkFlagsAsRequired(cmd, "bridgeAddress")
			}
			deployments = append(deployments, "erc721Handler")
		}
		if GenericHandler {
			if !Bridge {
				flags.MarkFlagsAsRequired(cmd, "bridgeAddress")
			}
			deployments = append(deployments, "genericHandler")
		}
		if Erc20 {
			flags.MarkFlagsAsRequired(cmd, "erc20Symbol", "erc20Name")
			deployments = append(deployments, "erc20")
		}
		if Erc721 {
			flags.MarkFlagsAsRequired(cmd, "erc721Name", "erc721Symbol", "erc721BaseURI")
			deployments = append(deployments, "erc721")
		}
	}

	if len(deployments) == 0 {
		log.Error().Err(ErrNoDeploymentFlagsProvided)
		return ErrNoDeploymentFlagsProvided
	}

	return nil
}

var deployments []string
var bridgeAddr common.Address
var relayerAddresses []common.Address

func ProcessDeployFlags(cmd *cobra.Command, args []string) error {
	if common.IsHexAddress(BridgeAddress) {
		bridgeAddr = common.HexToAddress(BridgeAddress)
	}
	for _, addr := range Relayers {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("invalid relayer address %s", addr)
		}
		relayerAddresses = append(relayerAddresses, common.HexToAddress(addr))
	}
	return nil
}

func CallDeployCLI(cmd *cobra.Command, args []string) error {
	txFabric := evmtransaction.NewTransaction
	return DeployCLI(cmd, args, txFabric, &evmgaspricer.LondonGasPriceDeterminant{})
}

func DeployCLI(cmd *cobra.Command, args []string, txFabric calls.TxFabric, gasPricer utils.GasPricerWithPostConfig) error {
	// fetch global flag values
	url, gasLimit, gasPrice, senderKeyPair, err := flags.GlobalFlagValues(cmd)
	if err != nil {
		return err
	}

	log.Debug().Msgf("url: %s gas limit: %v gas price: %v", url, gasLimit, gasPrice)
	log.Debug().Msgf("SENDER Private key 0x%s", hex.EncodeToString(crypto.FromECDSA(senderKeyPair.PrivateKey())))

	ethClient, err := evmclient.NewEVMClientFromParams(url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("ethereum client error: %v", err)).Msg("error initializing new EVM client")
		return err
	}
	gasPricer.SetClient(ethClient)
	gasPricer.SetOpts(&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice})
	log.Debug().Msgf("Relaysers for deploy %+v", Relayers)
	log.Debug().Msgf("all bool: %v", DeployAll)

	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			log.Debug().Msgf("deploying bridge..")
			bridgeAddr, err = calls.DeployBridge(ethClient, txFabric, gasPricer, DomainId, relayerAddresses, big.NewInt(0).SetUint64(RelayerThreshold), big.NewInt(0).SetUint64(Fee))
			if err != nil {
				log.Error().Err(fmt.Errorf("bridge deploy failed: %w", err))
				return err
			}
			deployedContracts["bridge"] = bridgeAddr.String()
			log.Debug().Msgf("bridge address; %v", bridgeAddr.String())
		case "erc20Handler":
			log.Debug().Msgf("deploying ERC20 handler..")
			erc20HandlerAddr, err := calls.DeployErc20Handler(ethClient, txFabric, gasPricer, bridgeAddr)
			if err != nil {
				log.Error().Err(fmt.Errorf("ERC20 handler deploy failed: %w", err))
				return err
			}
			deployedContracts["erc20Handler"] = erc20HandlerAddr.String()
		case "genericHandler":
			log.Debug().Msgf("deploying generic handler..")
			emptyAddr := common.Address{}
			if bridgeAddr == emptyAddr {
				log.Error().Err(errors.New("bridge flag or bridgeAddress param should be set for contracts deployments"))
				return err
			}

			genericHandlerAddr, err := calls.DeployGenericHandler(ethClient, txFabric, gasPricer, bridgeAddr)
			if err != nil {
				log.Error().Err(fmt.Errorf("Generic handler deploy failed: %w", err))
				return err
			}
			deployedContracts["genericHandler"] = genericHandlerAddr.String()
		case "erc20":
			log.Debug().Msgf("deploying ERC20..")
			erc20Addr, err := calls.DeployErc20(ethClient, txFabric, gasPricer, Erc20Name, Erc20Symbol)
			if err != nil {
				log.Error().Err(fmt.Errorf("erc 20 deploy failed: %w", err))
				return err
			}
			deployedContracts["erc20Token"] = erc20Addr.String()
		case "erc721":
			log.Debug().Msgf("deploying ERC721..")
			erc721Addr, err := calls.DeployErc721(ethClient, txFabric, gasPricer, Erc721Name, Erc721Symbol, Erc721BaseURI)
			if err != nil {
				log.Error().Err(fmt.Errorf("ERC721 deploy failed: %w", err))
				return err
			}

			deployedContracts["erc721Token"] = erc721Addr.String()
		case "erc721Handler":
			log.Debug().Msgf("deploying ERC721 handler..")
			erc721HandlerAddr, err := calls.DeployErc721Handler(ethClient, txFabric, gasPricer, bridgeAddr)
			if err != nil {
				log.Error().Err(fmt.Errorf("ERC721 handler deploy failed: %w", err))
				return err
			}
			deployedContracts["erc721Handler"] = erc721HandlerAddr.String()
		}

	}

	log.Info().Msgf("%+v", deployedContracts)
	return nil
}
