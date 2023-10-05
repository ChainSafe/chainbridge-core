package deposithandlers

import (
	"errors"
	"math/big"

	"github.com/ChainSafe/sygma-core/types"
)

type GenericDepositHandler struct {
	ArbitraryFunction arbitraryFunction
	Config            interface{}
}

// GenericDepositHandler converts data pulled from generic deposit event logs into message
func (dh *GenericDepositHandler) HandleDeposit(sourceID, destId uint8, nonce uint64, resourceID types.ResourceID, calldata, handlerResponse []byte) (*types.Message, error) {
	if len(calldata) < 32 {
		err := errors.New("invalid calldata length: less than 32 bytes")
		return nil, err
	}
	err := dh.ArbitraryFunction(dh.Config)
	if err != nil {
		return nil, err
	}
	// first 32 bytes are metadata length
	metadataLen := big.NewInt(0).SetBytes(calldata[:32])
	metadata := calldata[32 : 32+metadataLen.Int64()]
	payload := []interface{}{
		metadata,
	}

	// generic handler has specific payload length and doesn't support arbitrary metadata
	meta := types.Metadata{}
	return types.NewMessage(sourceID, destId, nonce, resourceID, types.GenericTransfer, payload, meta), nil
}
