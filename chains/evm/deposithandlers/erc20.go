package deposithandlers

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/sygma-core/types"
)

type Erc20DepositHandler struct {
	ArbitraryFunction arbitraryFunction
	Config            interface{}
}

// Erc20DepositHandler converts data pulled from event logs into message
// handlerResponse can be an empty slice
func (dh *Erc20DepositHandler) HandleDeposit(sourceID, destId uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*types.Message, error) {
	if len(calldata) < 84 {
		err := errors.New("invalid calldata length: less than 84 bytes")
		return nil, err
	}

	err := dh.ArbitraryFunction(dh.Config)
	if err != nil {
		return nil, err
	}

	// @dev
	// amount: first 32 bytes of calldata
	amount := calldata[:32]

	// lenRecipientAddress: second 32 bytes of calldata [32:64]
	// does not need to be derived because it is being calculated
	// within ERC20MessageHandler
	// https://github.com/ChainSafe/chainbridge-core/blob/main/chains/evm/voter/message-handler.go#L108

	// 32-64 is recipient address length
	recipientAddressLength := big.NewInt(0).SetBytes(calldata[32:64])

	// 64 - (64 + recipient address length) is recipient address
	recipientAddress := calldata[64:(64 + recipientAddressLength.Int64())]

	// if there is priority data, parse it and use it
	payload := []interface{}{
		amount,
		recipientAddress,
	}

	// arbitrary metadata that will be most likely be used by the relayer
	var metadata types.Metadata
	if 64+recipientAddressLength.Int64() < int64(len(calldata)) {
		priorityLength := big.NewInt(0).SetBytes(calldata[(64 + recipientAddressLength.Int64()):((64 + recipientAddressLength.Int64()) + 1)])

		// (64 + recipient address length + 1) - ((64 + recipient address length + 1) + priority length) is priority data
		priority := calldata[(64 + recipientAddressLength.Int64() + 1):((64 + recipientAddressLength.Int64()) + 1 + priorityLength.Int64())]

		// Assign the priority data to the Metadata struct
		metadata.Data = make(map[string]interface{})
		metadata.Data["Priority"] = priority[0]
	}
	return types.NewMessage(sourceID, destId, nonce, resourceID, types.FungibleTransfer, payload, metadata), nil
}
