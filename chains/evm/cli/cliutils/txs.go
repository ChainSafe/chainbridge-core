package cliutils

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

type DeployedContracts struct {
	BridgeAddress         common.Address
	ERC20HandlerAddress   common.Address
	ERC721HandlerAddress  common.Address
	GenericHandlerAddress common.Address
	ERC20TokenAddress     common.Address
}

// DeployContracts deploys Bridge, Relayer, ERC20Handler, ERC721Handler and CentrifugeAssetHandler and returns the addresses
func DeployContracts(client *evmclient.EVMClient, chainID uint8, initialRelayerThreshold *big.Int, relayerAddresses []common.Address, erc20Name, erc20Symbol string, bridgeFee *big.Int) (*DeployedContracts, error) {
	bridgeAddr, err := DeployBridge(client, chainID, relayerAddresses, initialRelayerThreshold, bridgeFee)
	if err != nil {
		return nil, err
	}

	erc20HandlerAddr, err := DeployERC20Handler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	erc721HandlerAddr, err := DeployERC721Handler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	genericHandlerAddr, err := DeployGenericHandler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	erc20Token, err := DeployERC20Token(client, erc20Name, erc20Symbol)
	if err != nil {
		return nil, err
	}

	dpc := &DeployedContracts{bridgeAddr, erc20HandlerAddr, erc721HandlerAddr, genericHandlerAddr, erc20Token}
	log.Debug().Msgf("Bridge %s \r\nerc20 handler %s \r\nerc721 handler %s \r\ngeneric handler %s \r\nerc20Contract %s", dpc.BridgeAddress.Hex(), dpc.ERC20HandlerAddress.Hex(), dpc.ERC721HandlerAddress.Hex(), dpc.GenericHandlerAddress.Hex(), dpc.ERC20TokenAddress.String())
	return dpc, nil
}

func DeployERC20Token(client *evmclient.EVMClient, name, symbol string) (common.Address, error) {
	log.Debug().Msgf("Deploying erc20..")
	bridgeAddr := common.Address{}
	return bridgeAddr, nil
}

func DeployBridge(client *evmclient.EVMClient, chainID uint8, relayerAddrs []common.Address, initialRelayerThreshold *big.Int, fee *big.Int) (common.Address, error) {
	log.Debug().Msgf("Deploying bridge..")
	bridgeAddr := common.Address{}
	return bridgeAddr, nil
}

func DeployERC20Handler(client *evmclient.EVMClient, bridgeAddress common.Address) (common.Address, error) {
	log.Debug().Msgf("Deploying erc20 handler..")
	erc20HandlerAddr := common.Address{}
	return erc20HandlerAddr, nil
}

func DeployERC721Handler(client *evmclient.EVMClient, bridgeAddress common.Address) (common.Address, error) {
	log.Debug().Msgf("Deploying erc721 handler..")
	erc721HandlerAddr := common.Address{}
	return erc721HandlerAddr, nil
}

func DeployERC721Token(client *evmclient.EVMClient) (common.Address, error) {
	log.Debug().Msgf("Deploying erc721..")
	erc721Addr := common.Address{}
	return erc721Addr, nil
}

func DeployGenericHandler(client *evmclient.EVMClient, bridgeAddress common.Address) (common.Address, error) {
	log.Debug().Msgf("Deploying generic handler..")
	genericHandlerAddr := common.Address{}
	return genericHandlerAddr, nil
}
