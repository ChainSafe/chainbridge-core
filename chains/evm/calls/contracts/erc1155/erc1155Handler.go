package erc1155

import (
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ERC1155HandlerContract struct {
	contracts.Contract
}

func NewERC1155HandlerContract(
	client calls.ContractCallerDispatcher,
	erc1155HandlerContractAddress common.Address,
	t transactor.Transactor,
) *ERC1155HandlerContract {
	a, _ := abi.JSON(strings.NewReader(consts.ERC1155HandlerABI))
	b := common.FromHex(consts.ERC1155HandlerBin)
	return &ERC1155HandlerContract{contracts.NewContract(erc1155HandlerContractAddress, a, b, client, t)}
}
