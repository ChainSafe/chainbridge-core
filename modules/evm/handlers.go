package evm

import (
	erc20Handler "github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/ERC20Handler"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func HandleErc20DepositedEvent(sourceID, destId uint8, nonce uint64, handlerContractAddress common.Address, backend bind.ContractBackend) (*DefaultEVMMessage, error) {
	contract, err := erc20Handler.NewERC20Handler(handlerContractAddress, backend)
	if err != nil {
		return nil, err
	}
	record, err := contract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		return nil, err
	}

	return &DefaultEVMMessage{
		Source:       sourceID,
		Destination:  destId,
		Type:         FungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   record.ResourceID,
		Payload: []interface{}{
			record.Amount.Bytes(),
			record.DestinationRecipientAddress,
		},
	}, nil
}
