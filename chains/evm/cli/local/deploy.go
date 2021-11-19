// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package local

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
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
	calls.ContractCallerClient
	evmgaspricer.GasPriceClient
	calls.ClientDeployer
}

func PrepareLocalEVME2EEnv(
	ethClient E2EClient,
	fabric calls.TxFabric,
	domainID uint8,
	treshHold *big.Int,
	mintTo common.Address,
) (EVME2EConfig, error) {
	bridgeAddr, err := deployBridgeForTest(ethClient, fabric, domainID, treshHold)
	if err != nil {
		return EVME2EConfig{}, err
	}

	erc20Addr, erc20HandlerAddr, err := PrepareErc20EVME2EEnv(ethClient, fabric, bridgeAddr, mintTo)
	if err != nil {
		return EVME2EConfig{}, err
	}

	erc721Addr, erc721HandlerAddr, err := PrepareErc721EVME2EEnv(ethClient, fabric, bridgeAddr, mintTo)
	if err != nil {
		return EVME2EConfig{}, err
	}

	assetStoreAddr, genericHandlerAddr, err := PrepareGenericEVME2EEnv(ethClient, fabric, bridgeAddr)
	if err != nil {
		return EVME2EConfig{}, err
	}

	log.Debug().Msgf("All deployments and preparations are done")

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

func PrepareErc20EVME2EEnv(ethClient E2EClient, fabric calls.TxFabric, bridgeAddr, mintTo common.Address) (common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	erc20Addr, erc20HandlerAddr, err := deployErc20ForTest(ethClient, fabric, bridgeAddr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	gasLimit := uint64(2000000)
	// Registering resource
	resourceID := calls.SliceTo32Bytes(append(common.LeftPadBytes(erc20Addr.Bytes(), 31), 0))
	registerResourceInput, err := calls.PrepareAdminSetResourceInput(erc20HandlerAddr, resourceID, erc20Addr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &bridgeAddr, registerResourceInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	// Minting tokens
	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	minInput, err := calls.PrepareMintTokensInput(mintTo, tenTokens)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &erc20Addr, minInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	// Approving tokens
	approveInput, err := calls.PrepareErc20ApproveInput(erc20HandlerAddr, tenTokens)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &erc20Addr, approveInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	// Adding minter
	minterInput, err := calls.PrepareErc20AddMinterInput(ethClient, erc20Addr, erc20HandlerAddr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &erc20Addr, minterInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	setBurnInput, err := calls.PrepareSetBurnableInput(erc20HandlerAddr, erc20Addr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &bridgeAddr, setBurnInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	return erc20Addr, erc20HandlerAddr, nil
}

func PrepareGenericEVME2EEnv(ethClient E2EClient, fabric calls.TxFabric, bridgeAddr common.Address) (common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	assetStoreAddr, genericHandlerAddr, err := deployGenericForTest(ethClient, fabric, bridgeAddr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	gasLimit := uint64(2000000)
	resourceID := calls.SliceTo32Bytes(append(common.LeftPadBytes(genericHandlerAddr.Bytes(), 31), 1))
	setGenericResourceInput, err := calls.PrepareAdminSetGenericResourceInput(
		genericHandlerAddr,
		resourceID,
		assetStoreAddr,
		[4]byte{0x65, 0x4c, 0xf8, 0x8c},
		big.NewInt(0),
		[4]byte{0x65, 0x4c, 0xf8, 0x8c},
	)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &bridgeAddr, setGenericResourceInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	return assetStoreAddr, genericHandlerAddr, nil
}

func PrepareErc721EVME2EEnv(ethClient E2EClient, fabric calls.TxFabric, bridgeAddr, mintTo common.Address) (common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(ethClient, nil)
	erc721Addr, erc721HandlerAddr, err := deployErc721ForTest(ethClient, fabric, bridgeAddr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	gasLimit := uint64(2000000)
	// Registering resource
	resourceID := calls.SliceTo32Bytes(append(common.LeftPadBytes(erc721Addr.Bytes(), 31), uint8(2)))
	registerResourceInput, err := calls.PrepareAdminSetResourceInput(erc721HandlerAddr, resourceID, erc721Addr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &bridgeAddr, registerResourceInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	// Adding minter
	_, err = calls.ERC721AddMinter(ethClient, fabric, staticGasPricer, gasLimit, erc721Addr, erc721HandlerAddr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	setBurnInput, err := calls.PrepareSetBurnableInput(erc721HandlerAddr, erc721Addr)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	_, err = calls.Transact(ethClient, fabric, staticGasPricer, &bridgeAddr, setBurnInput, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	return erc721Addr, erc721HandlerAddr, nil
}

func deployBridgeForTest(
	c E2EClient,
	fabric calls.TxFabric,
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
	fabric calls.TxFabric,
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
	fabric calls.TxFabric,
	bridgeAddr common.Address,
) (common.Address, common.Address, error) {
	staticGasPricer := evmgaspricer.NewStaticGasPriceDeterminant(c, nil)
	genericHandlerAddr, err := calls.DeployGenericHandler(c, fabric, staticGasPricer, bridgeAddr)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("Generic handler deploy failed: %w", err)
	}

	assetStoreAddr, err := calls.DeployCentrifugeAssetStore(c, fabric, staticGasPricer)
	if err != nil {
		return common.Address{}, common.Address{}, fmt.Errorf("Centrifuge asset store deploy failed: %w", err)
	}

	log.Debug().Msgf("Centrifuge asset store deployed to: %s; \n Generic Handler deployed to: %s", assetStoreAddr, genericHandlerAddr)
	return assetStoreAddr, genericHandlerAddr, nil
}

func deployErc721ForTest(
	c E2EClient,
	fabric calls.TxFabric,
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
