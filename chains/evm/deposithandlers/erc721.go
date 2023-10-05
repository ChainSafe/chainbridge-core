package deposithandlers

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/sygma-core/types"
)

type Erc721DepositHandler struct {
	ArbitraryFunction arbitraryFunction
	Config            interface{}
}

// Erc721DepositHandler converts data pulled from ERC721 deposit event logs into message
func (dh *Erc721DepositHandler) HandleDeposit(sourceID, destId uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*types.Message, error) {
	if len(calldata) < 64 {
		err := errors.New("invalid calldata length: less than 84 bytes")
		return nil, err
	}

	err := dh.ArbitraryFunction(dh.Config)
	if err != nil {
		return nil, err
	}

	// first 32 bytes are tokenId
	tokenId := calldata[:32]

	// 32 - 64 is recipient address length
	recipientAddressLength := big.NewInt(0).SetBytes(calldata[32:64])

	// 64 - (64 + recipient address length) is recipient address
	recipientAddress := calldata[64:(64 + recipientAddressLength.Int64())]

	// (64 + recipient address length) - ((64 + recipient address length) + 32) is metadata length
	metadataLength := big.NewInt(0).SetBytes(
		calldata[(64 + recipientAddressLength.Int64()):((64 + recipientAddressLength.Int64()) + 32)],
	)
	// ((64 + recipient address length) + 32) - ((64 + recipient address length) + 32 + metadata length) is metadata
	var metadata []byte
	var metadataStart int64
	if metadataLength.Cmp(big.NewInt(0)) == 1 {
		metadataStart = (64 + recipientAddressLength.Int64()) + 32
		metadata = calldata[metadataStart : metadataStart+metadataLength.Int64()]
	}
	// arbitrary metadata that will be most likely be used by the relayer
	var meta types.Metadata

	payload := []interface{}{
		tokenId,
		recipientAddress,
		metadata,
	}

	if 64+recipientAddressLength.Int64()+32+metadataLength.Int64() < int64(len(calldata)) {
		// (metadataStart + metadataLength) - (metadataStart + metadataLength + 1) is priority length
		priorityLength := big.NewInt(0).SetBytes(calldata[(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64()):(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64() + 1)])
		// (metadataStart + metadataLength + 1) - (metadataStart + metadataLength + 1) + priority length) is priority data
		priority := calldata[(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64() + 1):(64 + recipientAddressLength.Int64() + 32 + metadataLength.Int64() + 1 + priorityLength.Int64())]

		// Assign the priority data to the Metadata struct
		meta.Data = make(map[string]interface{})
		meta.Data["Priority"] = priority[0]
	}
	return types.NewMessage(sourceID, destId, nonce, resourceID, types.NonFungibleTransfer, payload, meta), nil
}
