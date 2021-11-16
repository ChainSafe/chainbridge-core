package calls

import (
	"context"
	"fmt"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func DeployCentrifugeAssetStore(c ClientDeployer, txFabric TxFabric, gasPriceClient GasPricer) (common.Address, error) {
	log.Debug().Msgf("Deploying Centrifuge asset store")
	parsed, err := abi.JSON(strings.NewReader(consts.CentrifugeAssetStoreABI))
	if err != nil {
		return common.Address{}, err
	}

	address, err := deployContract(c, parsed, common.FromHex(consts.CentrifugeAssetStoreBin), txFabric, gasPriceClient)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func prepareIsAssetStoredInput(hash [32]byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.CentrifugeAssetStoreABI))
	if err != nil {
		return []byte{}, err
	}

	input, err := a.Pack("_assetsStored", hash)
	if err != nil {
		return []byte{}, err
	}

	return input, nil
}

func parseIsAssetStoredOutput(output []byte) (bool, error) {
	a, err := abi.JSON(strings.NewReader(consts.CentrifugeAssetStoreABI))
	if err != nil {
		return false, err
	}

	res, err := a.Unpack("_assetsStored", output)
	if err != nil {
		return false, err
	}

	isAssetStored := *abi.ConvertType(res[0], new(bool)).(*bool)
	return isAssetStored, nil
}

func IsCentrifugeAssetStored(ethClient ContractCallerClient, storeAddr common.Address, hash [32]byte) (bool, error) {
	input, err := prepareIsAssetStoredInput(hash)
	if err != nil {
		log.Error().Err(fmt.Errorf("prepare input error: %v", err))
		return false, err
	}

	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &storeAddr,
		Data: input,
	}

	out, err := ethClient.CallContract(context.TODO(), ToCallArg(msg), nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("call contract error: %v", err))
		return false, err
	}

	isAssetStored, err := parseIsAssetStoredOutput(out)
	if err != nil {
		return false, nil
	}

	return isAssetStored, nil
}
