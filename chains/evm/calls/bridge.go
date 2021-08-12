package calls

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func PrepareSetBurnableInput(client ChainClient, handler, tokenAddress common.Address) ([]byte, error) {
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

func PrepareErc20DepositInput(destChainID uint8, resourceID [32]byte, data []byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("deposit", destChainID, resourceID, data)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}


func PrepareAddRelayerInput(relayer common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminAddRelayer", relayer)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}
