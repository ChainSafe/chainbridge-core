// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package local

import (
	"math/big"

	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/centrifuge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/erc721"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/generic"
	evmgaspricer "github.com/ChainSafe/sygma-core/chains/evm/calls/evmgaspricer"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor/signAndSend"
	"github.com/ChainSafe/sygma-core/keystore"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

var (
	DefaultRelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}
)

type BridgeConfig struct {
	BridgeAddr common.Address

	Erc20Addr        common.Address
	Erc20HandlerAddr common.Address
	Erc20ResourceID  types.ResourceID

	AssetStoreAddr     common.Address
	GenericHandlerAddr common.Address
	GenericResourceID  types.ResourceID

	Erc721Addr        common.Address
	Erc721HandlerAddr common.Address
	Erc721ResourceID  types.ResourceID

	ResourceIDERC721  string
	ResourceIDGeneric string
}

type EVMClient interface {
	calls.ContractCallerDispatcher
	evmgaspricer.GasPriceClient
}

func SetupEVMBridge(
	ethClient EVMClient,
	fabric calls.TxFabric,
	domainID uint8,
	threshold *big.Int,
	mintTo common.Address,
	relayerAddresses []common.Address,
) (BridgeConfig, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	t := signAndSend.NewSignAndSendTransactor(fabric, staticGasPricer, ethClient)

	bridgeContract := bridge.NewBridgeContract(ethClient, common.Address{}, t)
	bridgeContractAddress, err := bridgeContract.DeployContract(
		domainID, relayerAddresses, threshold, big.NewInt(0), big.NewInt(100),
	)
	if err != nil {
		return BridgeConfig{}, err
	}

	erc721Contract, erc721ContractAddress, erc721HandlerContractAddress, err := deployErc721(
		ethClient, t, bridgeContractAddress,
	)
	if err != nil {
		return BridgeConfig{}, err
	}

	erc20Contract, erc20ContractAddress, erc20HandlerContractAddress, err := deployErc20(
		ethClient, t, bridgeContractAddress,
	)

	if err != nil {
		return BridgeConfig{}, err
	}

	genericHandlerAddress, assetStoreAddress, err := deployGeneric(ethClient, t, bridgeContractAddress)
	if err != nil {
		return BridgeConfig{}, err
	}

	resourceIDERC20 := calls.SliceTo32Bytes(common.LeftPadBytes([]byte{0}, 31))
	resourceIDGenericHandler := calls.SliceTo32Bytes(common.LeftPadBytes([]byte{1}, 31))
	resourceIDERC721 := calls.SliceTo32Bytes(common.LeftPadBytes([]byte{2}, 31))

	conf := BridgeConfig{
		BridgeAddr: bridgeContractAddress,

		Erc20Addr:        erc20ContractAddress,
		Erc20HandlerAddr: erc20HandlerContractAddress,
		Erc20ResourceID:  resourceIDERC20,

		GenericHandlerAddr: genericHandlerAddress,
		AssetStoreAddr:     assetStoreAddress,
		GenericResourceID:  resourceIDGenericHandler,

		Erc721Addr:        erc721ContractAddress,
		Erc721HandlerAddr: erc721HandlerContractAddress,
		Erc721ResourceID:  resourceIDERC721,
	}

	err = SetupERC20Handler(bridgeContract, erc20Contract, mintTo, conf, resourceIDERC20)
	if err != nil {
		return BridgeConfig{}, err
	}

	err = SetupERC721Handler(bridgeContract, erc721Contract, conf, resourceIDERC721)
	if err != nil {
		return BridgeConfig{}, err
	}

	err = SetupGenericHandler(bridgeContract, conf, resourceIDGenericHandler)
	if err != nil {
		return BridgeConfig{}, err
	}

	_, err = bridgeContract.AdminChangeRelayerThreshold(threshold.Uint64(), transactor.TransactOptions{})
	if err != nil {
		return BridgeConfig{}, err
	}

	log.Debug().Msgf("All deployments and preparations are done")
	return conf, nil
}

func deployGeneric(
	ethClient EVMClient, t transactor.Transactor, bridgeContractAddress common.Address,
) (common.Address, common.Address, error) {
	assetStoreContract := centrifuge.NewAssetStoreContract(ethClient, common.Address{}, t)
	assetStoreAddress, err := assetStoreContract.DeployContract()
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	genericHandlerContract := generic.NewGenericHandlerContract(ethClient, common.Address{}, t)
	genericHandlerAddress, err := genericHandlerContract.DeployContract(bridgeContractAddress)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	log.Debug().Msgf(
		"Centrifuge asset store deployed to: %s; \n Generic Handler deployed to: %s",
		assetStoreAddress, genericHandlerAddress,
	)
	return genericHandlerAddress, assetStoreAddress, nil
}

func deployErc20(
	ethClient EVMClient, t transactor.Transactor, bridgeContractAddress common.Address,
) (*erc20.ERC20Contract, common.Address, common.Address, error) {
	erc20Contract := erc20.NewERC20Contract(ethClient, common.Address{}, t)
	erc20ContractAddress, err := erc20Contract.DeployContract("Test", "TST")
	if err != nil {
		return nil, common.Address{}, common.Address{}, err
	}
	erc20HandlerContract := erc20.NewERC20HandlerContract(ethClient, common.Address{}, t)
	erc20HandlerContractAddress, err := erc20HandlerContract.DeployContract(bridgeContractAddress)
	if err != nil {
		return nil, common.Address{}, common.Address{}, err
	}
	log.Debug().Msgf(
		"Erc20 deployed to: %s; \n Erc20 Handler deployed to: %s",
		erc20ContractAddress, erc20HandlerContractAddress,
	)
	return erc20Contract, erc20ContractAddress, erc20HandlerContractAddress, nil
}

func deployErc721(
	ethClient EVMClient, t transactor.Transactor, bridgeContractAddress common.Address,
) (*erc721.ERC721Contract, common.Address, common.Address, error) {
	erc721Contract := erc721.NewErc721Contract(ethClient, common.Address{}, t)
	erc721ContractAddress, err := erc721Contract.DeployContract("TestERC721", "TST721", "")
	if err != nil {
		return nil, common.Address{}, common.Address{}, err
	}
	erc721HandlerContract := erc721.NewERC721HandlerContract(ethClient, common.Address{}, t)
	erc721HandlerContractAddress, err := erc721HandlerContract.DeployContract(bridgeContractAddress)
	if err != nil {
		return nil, common.Address{}, common.Address{}, err
	}
	log.Debug().Msgf(
		"Erc721 deployed to: %s; \n Erc721 Handler deployed to: %s",
		erc721ContractAddress, erc721HandlerContractAddress,
	)
	return erc721Contract, erc721ContractAddress, erc721HandlerContractAddress, nil
}

func SetupERC20Handler(
	bridgeContract *bridge.BridgeContract, erc20Contract *erc20.ERC20Contract, mintTo common.Address, conf BridgeConfig, resourceID types.ResourceID,
) error {
	_, err := bridgeContract.AdminSetResource(
		conf.Erc20HandlerAddr, resourceID, conf.Erc20Addr, transactor.TransactOptions{GasLimit: 2000000},
	)
	if err != nil {
		return err
	}
	// Minting tokens
	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	_, err = erc20Contract.MintTokens(mintTo, tenTokens, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	// Approving tokens
	_, err = erc20Contract.ApproveTokens(conf.Erc20HandlerAddr, tenTokens, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	// Adding minter
	_, err = erc20Contract.AddMinter(conf.Erc20HandlerAddr, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	// Set burnable input
	_, err = bridgeContract.SetBurnableInput(conf.Erc20HandlerAddr, conf.Erc20Addr, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	return nil
}

func SetupGenericHandler(bridgeContract *bridge.BridgeContract, conf BridgeConfig, resourceID types.ResourceID) error {
	_, err := bridgeContract.AdminSetGenericResource(
		conf.GenericHandlerAddr,
		resourceID,
		conf.AssetStoreAddr,
		[4]byte{0x65, 0x4c, 0xf8, 0x8c},
		big.NewInt(0),
		[4]byte{0x65, 0x4c, 0xf8, 0x8c},
		transactor.TransactOptions{GasLimit: 2000000},
	)
	if err != nil {
		return err
	}
	return nil
}

func SetupERC721Handler(bridgeContract *bridge.BridgeContract, erc721Contract *erc721.ERC721Contract, conf BridgeConfig, resourceID types.ResourceID) error {
	_, err := bridgeContract.AdminSetResource(conf.Erc721HandlerAddr, resourceID, conf.Erc721Addr, transactor.TransactOptions{GasLimit: 2000000})
	if err != nil {
		return err
	}
	// Adding minter
	_, err = erc721Contract.AddMinter(conf.Erc721HandlerAddr, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	// Set burnable input
	_, err = bridgeContract.SetBurnableInput(conf.Erc721HandlerAddr, conf.Erc721Addr, transactor.TransactOptions{})
	if err != nil {
		return err
	}
	return nil
}
