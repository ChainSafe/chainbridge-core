package calls

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

// @dev
// inputs here and in bridge.go could get consolidated into something similar to txFabric in deploy.go

func PrepareMintTokensInput(destAddr common.Address, amount *big.Int) ([]byte, error) {
	log.Debug().Msgf("Minting tokens %s %s", destAddr.String(), amount.String())
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("mint", destAddr, amount)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareErc20ApproveInput(target common.Address, amount *big.Int) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("approve", target, amount)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareErc20AddMinterInput(client ChainClient, erc20Contract, handler common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	role, err := MinterRole(client, erc20Contract)
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("grantRole", role, handler)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareErc20DepositInput(bridgeAddress, recipientAddress common.Address, amount *big.Int, rId [32]byte, destChainId uint8) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("deposit", bridgeAddress, recipientAddress, amount, rId, destChainId)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareRegisterGenericResourceInput(handler common.Address, rId [32]byte, addr common.Address, depositSig, executeSig [4]byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err // Not sure what status to use here
	}
	input, err := a.Pack("adminSetGenericResource", handler, rId, addr, depositSig, executeSig)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareERC20BalanceInput(erc20Addr, accountAddr common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("balanceOf", erc20Addr, accountAddr)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func MinterRole(chainClient ChainClient, erc20Contract common.Address) ([32]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return [32]byte{}, err
	}
	input, err := a.Pack("MINTER_ROLE")
	if err != nil {
		return [32]byte{}, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &erc20Contract, Data: input}
	out, err := chainClient.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return [32]byte{}, err
	}
	res, err := a.Unpack("MINTER_ROLE", out)
	if err != nil {
		return [32]byte{}, err
	}
	out0 := *abi.ConvertType(res[0], new([32]byte)).(*[32]byte)
	return out0, nil
}