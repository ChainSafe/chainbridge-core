// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package local

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contract"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"math/big"
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

type EVME2EConfig struct {
	BridgeAddr         common.Address
	Erc20Addr          common.Address
	Erc20HandlerAddr   common.Address
	AssetStoreAddr     common.Address
	GenericHandlerAddr common.Address
	Erc721Addr         common.Address
	Erc721HandlerAddr  common.Address
}

type E2EClient interface {
	client.ContractCallerClient
	evmgaspricer.GasPriceClient
	client.ClientDeployer
}

func PrepareLocalEVME2EEnv(
	ethClient E2EClient,
	fabric client.TxFabric,
	domainID uint8,
	threshold *big.Int,
	mintTo common.Address,
) (EVME2EConfig, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	t := transactor.NewSignAndSendTransactor(fabric, staticGasPricer, ethClient)

	bridgeContract := bridge.NewBridgeContract(ethClient, common.Address{}, t)
	bridgeContractAddress, err := bridgeContract.DeployContract(
		domainID, DefaultRelayerAddresses, threshold, big.NewInt(0), big.NewInt(100),
	)
	if err != nil {
		return EVME2EConfig{}, err
	}

	erc721Contract, erc721ContractAddress, erc721HandlerContractAddress, err := deployErc721(ethClient, t)
	if err != nil {
		return EVME2EConfig{}, err
	}

	erc20Contract, erc20ContractAddress, erc20HandlerContractAddress, err := deployErc20(ethClient, t)
	if err != nil {
		return EVME2EConfig{}, err
	}

	genericHandlerAddress, assetStoreAddress, err := deployGeneric(ethClient, t, bridgeContractAddress)
	if err != nil {
		return EVME2EConfig{}, err
	}

	conf := EVME2EConfig{
		BridgeAddr: bridgeContractAddress,

		Erc20Addr:        erc20ContractAddress,
		Erc20HandlerAddr: erc20HandlerContractAddress,

		GenericHandlerAddr: genericHandlerAddress,
		AssetStoreAddr:     assetStoreAddress,

		Erc721Addr:        erc721ContractAddress,
		Erc721HandlerAddr: erc721HandlerContractAddress,
	}

	err = PrepareErc20EVME2EEnv(bridgeContract, erc20Contract, mintTo, conf)
	if err != nil {
		return EVME2EConfig{}, err
	}

	err = PrepareErc721EVME2EEnv(bridgeContract, erc721Contract, conf)
	if err != nil {
		return EVME2EConfig{}, err
	}

	err = PrepareGenericEVME2EEnv(bridgeContract, conf)
	if err != nil {
		return EVME2EConfig{}, err
	}

	log.Debug().Msgf("All deployments and preparations are done")
	return conf, nil
}

func deployGeneric(
	ethClient E2EClient, t transactor.Transactor, bridgeContractAddress common.Address,
) (common.Address, common.Address, error) {
	genericHandlerAddress, err := contract.DeployContract(
		consts.GenericHandlerABI, consts.GenericHandlerBin, ethClient, t, bridgeContractAddress,
	)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	assetStoreAddress, err := contract.DeployContract(
		consts.CentrifugeAssetStoreABI, consts.CentrifugeAssetStoreBin, ethClient, t, bridgeContractAddress,
	)
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
	ethClient E2EClient, t transactor.Transactor,
) (*erc20.ERC20Contract, common.Address, common.Address, error) {
	erc20Contract := erc20.NewERC20Contract(ethClient, common.Address{}, t)
	erc20ContractAddress, err := erc20Contract.DeployContract("Test", "TST")
	if err != nil {
		return nil, common.Address{}, common.Address{}, err
	}
	erc20HandlerContractAddress, err := contract.DeployContract(
		consts.ERC20HandlerABI, consts.ERC20HandlerBin, ethClient, t, erc20ContractAddress,
	)
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
	ethClient E2EClient, t transactor.Transactor,
) (*erc721.ERC721Contract, common.Address, common.Address, error) {
	erc721Contract := erc721.NewErc721Contract(ethClient, common.Address{}, t)
	erc721ContractAddress, err := erc721Contract.DeployContract("TestERC721", "TST721", "")
	if err != nil {
		return nil, common.Address{}, common.Address{}, err
	}
	erc721HandlerContractAddress, err := contract.DeployContract(
		consts.HandlerABI, consts.HandlerBin, ethClient, t, erc721ContractAddress,
	)
	if err != nil {
		return nil, common.Address{}, common.Address{}, err
	}
	log.Debug().Msgf(
		"Erc721 deployed to: %s; \n Erc721 Handler deployed to: %s",
		erc721ContractAddress, erc721HandlerContractAddress,
	)
	return erc721Contract, erc721ContractAddress, erc721HandlerContractAddress, nil
}

func PrepareErc20EVME2EEnv(
	bridgeContract *bridge.BridgeContract, erc20Contract *erc20.ERC20Contract, mintTo common.Address, conf EVME2EConfig,
) error {
	// Setting resource
	resourceID := client.SliceTo32Bytes(append(common.LeftPadBytes(conf.Erc20Addr.Bytes(), 31), 0))
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

func PrepareGenericEVME2EEnv(bridgeContract *bridge.BridgeContract, conf EVME2EConfig) error {
	resourceID := client.SliceTo32Bytes(append(common.LeftPadBytes(conf.GenericHandlerAddr.Bytes(), 31), 1))
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

func PrepareErc721EVME2EEnv(bridgeContract *bridge.BridgeContract, erc721Contract *erc721.ERC721Contract, conf EVME2EConfig) error {
	// Registering resource
	resourceID := client.SliceTo32Bytes(append(common.LeftPadBytes(conf.Erc20Addr.Bytes(), 31), uint8(2)))
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
