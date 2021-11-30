// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package local

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/centrifuge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc20"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/keystore"
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
	conf, err := DeployLocalEVME2EEnv(ethClient, fabric, domainID, threshold)
	if err != nil {
		return EVME2EConfig{}, err
	}

	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	t := transactor.NewSignAndSendTransactor(fabric, staticGasPricer, ethClient)

	bridgeContract := bridge.NewBridgeContract(ethClient, conf.BridgeAddr, t)
	erc721Contract := erc721.NewErc721Contract(ethClient, conf.Erc721Addr, t)

	err = PrepareErc20EVME2EEnv(ethClient, fabric, mintTo, conf)
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

func DeployLocalEVME2EEnv(
	ethClient E2EClient,
	fabric client.TxFabric,
	domainID uint8,
	threshold *big.Int,
) (EVME2EConfig, error) {
	bridgeAddr, err := deployBridgeForTest(ethClient, fabric, domainID, threshold)
	if err != nil {
		return EVME2EConfig{}, err
	}

	erc20Addr, erc20HandlerAddr, err := deployErc20ForTest(ethClient, fabric, bridgeAddr)
	if err != nil {
		return EVME2EConfig{}, err
	}

	erc721Addr, erc721HandlerAddr, err := deployErc721ForTest(ethClient, fabric, bridgeAddr)
	if err != nil {
		return EVME2EConfig{}, err
	}

	assetStoreAddr, genericHandlerAddr, err := deployGenericForTest(ethClient, fabric, bridgeAddr)
	if err != nil {
		return EVME2EConfig{}, err
	}

	return EVME2EConfig{
		BridgeAddr: bridgeAddr,

		Erc20Addr:        erc20Addr,
		Erc20HandlerAddr: erc20HandlerAddr,

		GenericHandlerAddr: genericHandlerAddr,
		AssetStoreAddr:     assetStoreAddr,

		Erc721Addr:        erc721Addr,
		Erc721HandlerAddr: erc721HandlerAddr,
	}, nil
}

func PrepareErc20EVME2EEnv(ethClient E2EClient, fabric client.TxFabric, mintTo common.Address, conf EVME2EConfig) error {
	// TODO - will be moved once Erc20Contract is refactored
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	t := transactor.NewSignAndSendTransactor(fabric, staticGasPricer, ethClient)
	bridgeContract := bridge.NewBridgeContract(ethClient, conf.BridgeAddr, t)

	// Setting resource
	resourceID := client.SliceTo32Bytes(append(common.LeftPadBytes(conf.Erc20Addr.Bytes(), 31), 0))
	_, err := bridgeContract.AdminSetResource(conf.Erc20HandlerAddr, resourceID, conf.Erc20Addr, transactor.TransactOptions{GasLimit: 2000000})
	if err != nil {
		return err
	}

	// Minting tokens
	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	minInput, err := erc20.PrepareMintTokensInput(mintTo, tenTokens)
	if err != nil {
		return err
	}
	_, err = client.Transact(ethClient, fabric, staticGasPricer, &conf.Erc20Addr, minInput, 2000000, big.NewInt(0))
	if err != nil {
		return err
	}

	// Approving tokens
	approveInput, err := erc20.PrepareErc20ApproveInput(conf.Erc20HandlerAddr, tenTokens)
	if err != nil {
		return err
	}
	_, err = client.Transact(ethClient, fabric, staticGasPricer, &conf.Erc20Addr, approveInput, 2000000, big.NewInt(0))
	if err != nil {
		return err
	}

	// Adding minter
	minterInput, err := erc20.PrepareErc20AddMinterInput(ethClient, conf.Erc20Addr, conf.Erc20HandlerAddr)
	if err != nil {
		return err
	}
	_, err = client.Transact(ethClient, fabric, staticGasPricer, &conf.Erc20Addr, minterInput, 2000000, big.NewInt(0))
	if err != nil {
		return err
	}

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

func deployBridgeForTest(
	c E2EClient,
	fabric client.TxFabric,
	domainID uint8,
	treshHold *big.Int,
) (common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(c, nil)
	bridgeAddr, err := calls.DeployBridge(c, fabric, staticGasPricer, domainID, DefaultRelayerAddresses, treshHold, big.NewInt(0))
	if err != nil {
		return common.Address{}, fmt.Errorf("Bridge deploy failed: %w", err)
	}

	log.Debug().Msgf("Bridge deployed to address: %s", bridgeAddr)
	return bridgeAddr, nil
}

func deployErc20ForTest(
	c E2EClient,
	fabric client.TxFabric,
	bridgeAddr common.Address,
) (common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(c, nil)
	erc20Addr, err := calls.DeployErc20(c, fabric, staticGasPricer, "Test", "TST")
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("ERC20 deploy failed: %w", err)
	}

	erc20HandlerAddr, err := calls.DeployErc20Handler(c, fabric, staticGasPricer, bridgeAddr)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("ERC20 handler deploy failed: %w", err)
	}

	log.Debug().Msgf("Erc20 deployed to: %s; \n Erc20 Handler deployed to: %s", erc20Addr, erc20HandlerAddr)
	return erc20Addr, erc20HandlerAddr, nil
}

func deployGenericForTest(
	c E2EClient,
	fabric client.TxFabric,
	bridgeAddr common.Address,
) (common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(c, nil)
	genericHandlerAddr, err := calls.DeployGenericHandler(c, fabric, staticGasPricer, bridgeAddr)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("Generic handler deploy failed: %w", err)
	}

	assetStoreAddr, err := centrifuge.DeployCentrifugeAssetStore(c, fabric, staticGasPricer)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("Centrifuge asset store deploy failed: %w", err)
	}

	log.Debug().Msgf("Centrifuge asset store deployed to: %s; \n Generic Handler deployed to: %s", assetStoreAddr, genericHandlerAddr)
	return assetStoreAddr, genericHandlerAddr, nil
}

func deployErc721ForTest(
	c E2EClient,
	fabric client.TxFabric,
	bridgeAddr common.Address,
) (common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(c, nil)
	erc721Addr, err := calls.DeployErc721(c, fabric, staticGasPricer, "TestERC721", "TST721", "")
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("ERC721 deploy failed: %w", err)
	}

	erc721HandlerAddr, err := calls.DeployErc721Handler(c, fabric, staticGasPricer, bridgeAddr)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("ERC721 handler deploy failed: %w", err)
	}

	log.Debug().Msgf("Erc721 deployed to: %s; \n Erc721 Handler deployed to: %s", erc721Addr, erc721HandlerAddr)
	return erc721Addr, erc721HandlerAddr, nil
}
