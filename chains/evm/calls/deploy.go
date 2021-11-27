package calls

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

func DeployErc20(c ClientDeployer, txFabric TxFabric, gasPriceClient GasPricer, name, symbol string) (common.Address, error) {
	parsed, err := abi.JSON(strings.NewReader(consts.ERC20PresetMinterPauserABI))
	if err != nil {
		return common.Address{}, err
	}
	address, err := deployContract(c, parsed, common.FromHex(consts.ERC20PresetMinterPauserBin), txFabric, gasPriceClient, name, symbol)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func DeployBridge(c ClientDeployer, txFabric TxFabric, gasPriceClient GasPricer, domainID uint8, relayerAddrs []common.Address, initialRelayerThreshold *big.Int, fee *big.Int) (common.Address, error) {
	parsed, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return common.Address{}, err
	}
	address, err := deployContract(c, parsed, common.FromHex(consts.BridgeBin), txFabric, gasPriceClient, domainID, relayerAddrs, initialRelayerThreshold, fee, big.NewInt(100))
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func DeployErc20Handler(c ClientDeployer, txFabric TxFabric, gasPriceClient GasPricer, bridgeAddress common.Address) (common.Address, error) {
	log.Debug().Msgf("Deploying ERC20 Handler with params: %s", bridgeAddress.String())
	parsed, err := abi.JSON(strings.NewReader(consts.ERC20HandlerABI))
	if err != nil {
		return common.Address{}, err
	}
	address, err := deployContract(c, parsed, common.FromHex(consts.ERC20HandlerBin), txFabric, gasPriceClient, bridgeAddress)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func DeployErc721(c ClientDeployer, txFabric TxFabric, gasPriceClient GasPricer, name, symbol, baseURI string) (common.Address, error) {
	parsed, err := abi.JSON(strings.NewReader(consts.ERC721PresetMinterPauserABI))
	if err != nil {
		return common.Address{}, err
	}
	address, err := deployContract(c, parsed, common.FromHex(consts.ERC721PresetMinterPauserBin), txFabric, gasPriceClient, name, symbol, baseURI)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func DeployErc721Handler(c ClientDeployer, txFabric TxFabric, gasPriceClient GasPricer, bridgeAddress common.Address) (common.Address, error) {
	log.Debug().Msgf("Deploying ERC721 Handler with params: %s", bridgeAddress.String())
	parsed, err := abi.JSON(strings.NewReader(consts.ERC721HandlerABI))
	if err != nil {
		return common.Address{}, err
	}
	address, err := deployContract(c, parsed, common.FromHex(consts.ERC721HandlerBin), txFabric, gasPriceClient, bridgeAddress)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func DeployGenericHandler(c ClientDeployer, txFabric TxFabric, gasPriceClient GasPricer, bridgeAddress common.Address) (common.Address, error) {
	log.Debug().Msgf("Deploying Generic Handler with params: %s", bridgeAddress.String())
	parsed, err := abi.JSON(strings.NewReader(consts.GenericHandlerABI))
	if err != nil {
		return common.Address{}, err
	}
	address, err := deployContract(c, parsed, common.FromHex(consts.GenericHandlerBin), txFabric, gasPriceClient, bridgeAddress)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func deployContract(client ClientDeployer, abi abi.ABI, bytecode []byte, txFabric TxFabric, gasPriceClient GasPricer, params ...interface{}) (common.Address, error) {
	defer client.UnlockNonce()

	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return common.Address{}, err
	}
	input, err := abi.Pack("", params...)
	if err != nil {
		return common.Address{}, err
	}

	gp, err := gasPriceClient.GasPrice()
	if err != nil {
		return common.Address{}, err
	}
	tx, err := txFabric(n.Uint64(), nil, big.NewInt(0), consts.DefaultDeployGasLimit, gp, append(bytecode, input...))
	if err != nil {
		return common.Address{}, err
	}
	hash, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return common.Address{}, err
	}
	time.Sleep(2 * time.Second)
	_, err = client.WaitAndReturnTxReceipt(tx.Hash())
	if err != nil {
		return common.Address{}, err
	}
	log.Debug().Str("hash", hash.String()).Uint64("nonce", n.Uint64()).Msgf("Contract deployed")
	address := crypto.CreateAddress(client.From(), n.Uint64())
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return common.Address{}, err
	}
	// checks bytecode at address
	// nil is latest block
	if code, err := client.CodeAt(context.Background(), address, nil); err != nil {
		return common.Address{}, err
	} else if len(code) == 0 {
		return common.Address{}, fmt.Errorf("no code at provided address %s", address.String())
	}
	return address, nil
}
