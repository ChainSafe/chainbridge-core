package generic

import (
	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/consts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

type GenericHandlerContract struct {
	contracts.Contract
}

func NewGenericHandlerContract(
	client calls.ContractCallerDispatcher,
	assetStoreContractAddress common.Address,
	transactor transactor.Transactor,
) *GenericHandlerContract {
	a, _ := abi.JSON(strings.NewReader(consts.GenericHandlerABI))
	b := common.FromHex(consts.GenericHandlerBin)
	return &GenericHandlerContract{contracts.NewContract(assetStoreContractAddress, a, b, client, transactor)}
}
