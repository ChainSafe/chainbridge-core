package calls

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
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
	log.Debug().Msgf("ResourceID %x", rId)
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

func PrepareIsRelayerInput(address common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}

	data, err := a.Pack("isRelayer", address)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return nil, err
	}
	return data, nil
}

func ParseIsRelayerOutput(output []byte) (bool, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return false, err
	}

	res, err := a.Unpack("isRelayer", output)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return false, err
	}

	b := abi.ConvertType(res[0], new(bool)).(*bool)
	return *b, nil
}

func PrepareSetDepositNonceInput(domainID uint8, depositNonce uint64) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetDepositNonce", domainID, depositNonce)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}
