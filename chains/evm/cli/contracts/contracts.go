package contracts

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/centrifuge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"math/big"
)

func InitializeErc20Contract(
	url string,
	gasLimit uint64,
	gasPrice *big.Int,
	senderKeyPair *secp256k1.Keypair,
	erc20ContractAddress common.Address,
) (*erc20.ERC20Contract, error) {
	ethClient, trans, err := InitializeClientAndTransactor(url, gasLimit, gasPrice, senderKeyPair)
	if err != nil {
		return nil, err
	}
	erc20Contract := erc20.NewERC20Contract(ethClient, erc20ContractAddress, trans)

	return erc20Contract, nil
}

func InitializeErc721Contract(
	url string,
	gasLimit uint64,
	gasPrice *big.Int,
	senderKeyPair *secp256k1.Keypair,
	erc721ContractAddress common.Address,
) (*erc721.ERC721Contract, error) {
	ethClient, trans, err := InitializeClientAndTransactor(url, gasLimit, gasPrice, senderKeyPair)
	if err != nil {
		return nil, err
	}
	erc721Contract := erc721.NewErc721Contract(ethClient, erc721ContractAddress, trans)

	return erc721Contract, nil
}

func InitializeAssetStoreContract(
	url string,
	gasLimit uint64,
	gasPrice *big.Int,
	senderKeyPair *secp256k1.Keypair,
	assetStoreContractAddress common.Address,
) (*centrifuge.AssetStoreContract, error) {
	ethClient, trans, err := InitializeClientAndTransactor(url, gasLimit, gasPrice, senderKeyPair)
	if err != nil {
		return nil, err
	}

	assetStoreContract := centrifuge.NewAssetStoreContract(ethClient, assetStoreContractAddress, trans)
	return assetStoreContract, nil
}

func InitializeBridgeContract(
	url string,
	gasLimit uint64,
	gasPrice *big.Int,
	senderKeyPair *secp256k1.Keypair,
	bridgeContractAddress common.Address,
) (*bridge.BridgeContract, error) {
	ethClient, trans, err := InitializeClientAndTransactor(url, gasLimit, gasPrice, senderKeyPair)
	if err != nil {
		return nil, err
	}
	bridgeContract := bridge.NewBridgeContract(ethClient, bridgeContractAddress, trans)
	return bridgeContract, nil
}

func InitializeClientAndTransactor(
	url string,
	gasLimit uint64,
	gasPrice *big.Int,
	senderKeyPair *secp256k1.Keypair,
) (*evmclient.EVMClient, transactor.Transactor, error) {
	txFabric := evmtransaction.NewTransaction

	ethClient, err := evmclient.NewEVMClientFromParams(
		url, senderKeyPair.PrivateKey())
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client intialization error: %v", err))
		return nil, nil, err
	}

	gasPricer := evmgaspricer.NewLondonGasPriceClient(
		ethClient,
		&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice},
	)

	trans := transactor.NewSignAndSendTransactor(txFabric, gasPricer, ethClient)
	return ethClient, trans, nil
}
