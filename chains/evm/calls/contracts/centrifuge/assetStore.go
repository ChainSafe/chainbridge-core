package centrifuge

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type AssetStoreContract struct {
	contracts.Contract
}

func NewAssetStoreContract(
	client client.ContractCallerDispatcherClient,
	assetStoreContractAddress common.Address,
	transactor transactor.Transactor,
) *AssetStoreContract {
	a, err := abi.JSON(strings.NewReader(consts.CentrifugeAssetStoreABI))
	if err != nil {
		log.Fatal().Msg("Unable to load AssetStore ABI") // TODO
	}
	b := common.FromHex(consts.CentrifugeAssetStoreBin)
	return &AssetStoreContract{contracts.NewContract(assetStoreContractAddress, a, b, client, transactor)}
}

func (c AssetStoreContract) IsCentrifugeAssetStored(hash [32]byte) (bool, error) {
	res, err := c.CallContract("_assetsStored", hash)
	if err != nil {
		return false, err
	}

	isAssetStored := *abi.ConvertType(res[0], new(bool)).(*bool)
	return isAssetStored, nil
}
