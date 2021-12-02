package contract

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

func DeployContract(
	abiString string,
	bytecodeString string,
	client client.ContractCallerDispatcherClient,
	t transactor.Transactor,
	params ...interface{},
) (common.Address, error) {
	a, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return common.Address{}, err
	}
	bytecode := common.FromHex(bytecodeString)
	contract := NewContract(common.Address{}, a, bytecode, client, t)
	contractAddress, err := contract.DeployContract(params)
	if err != nil {
		return common.Address{}, err
	}
	return contractAddress, nil
}
