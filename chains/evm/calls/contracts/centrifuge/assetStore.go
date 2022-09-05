package centrifuge

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type AssetStoreContract struct {
	contracts.Contract
}

func NewAssetStoreContract(
	client calls.ContractCallerDispatcher,
	assetStoreContractAddress common.Address,
	transactor transactor.Transactor,
) *AssetStoreContract {
	a, _ := abi.JSON(strings.NewReader(consts.CentrifugeAssetStoreABI))
	b := common.FromHex(consts.CentrifugeAssetStoreBin)
	return &AssetStoreContract{contracts.NewContract(assetStoreContractAddress, a, b, client, transactor)}
}

func (c *AssetStoreContract) IsCentrifugeAssetStored(hash [32]byte) (bool, error) {
	log.Debug().
		Str("hash", hexutil.Encode(hash[:])).
		Msgf("Getting is centrifuge asset stored")
	res, err := c.CallContract("_assetsStored", hash)
	if err != nil {
		return false, err
	}

	isAssetStored := *abi.ConvertType(res[0], new(bool)).(*bool)
	return isAssetStored, nil
}
