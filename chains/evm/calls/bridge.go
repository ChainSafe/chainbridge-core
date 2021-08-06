package calls

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// @dev
// inputs here and in erc20.go could get consolidated into something similar to txFabric in deploy.go

func PrepareSetBurnableInput(client ChainClient, bridge, handler, tokenAddress common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetBurnable", handler, tokenAddress)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareAdminSetResourceInput(handler common.Address, rId [32]byte, addr common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetResource", handler, rId, addr)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareErc20DepositInput(bridgeAddress, recipientAddress common.Address, amount *big.Int, rId [32]byte, destChainId uint8) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("deposit", bridgeAddress, recipientAddress, amount, rId, destChainId)
	if err != nil {
		return []byte{}, err
	}
	input = append(input, common.FromHex(ERC20PresetMinterPauserBin)...)
	return input, nil
}
