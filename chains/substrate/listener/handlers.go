package listener

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains/substrate"
	"github.com/ChainSafe/chainbridge-core/relayer"
)

func FungibleTransferHandler(sourceID uint8, evtI interface{}) (*relayer.Message, error) {
	evt, ok := evtI.(substrate.EventFungibleTransfer)
	if !ok {
		return nil, fmt.Errorf("failed to cast EventFungibleTransfer type")
	}
	//recipient := []byte{evt.Recipient[:]}
	return &relayer.Message{
		Source:       sourceID,
		Destination:  uint8(evt.Destination),
		DepositNonce: uint64(evt.DepositNonce),
		ResourceId:   evt.ResourceId,
		Payload: []interface{}{
			evt.Amount.Bytes(),
			[]byte(evt.Recipient),
		},
	}, nil
}

func NonFungibleTransferHandler(sourceID uint8, evtI interface{}) (*relayer.Message, error) {
	evt, ok := evtI.(substrate.EventNonFungibleTransfer)
	if !ok {
		return nil, fmt.Errorf("failed to cast EventNonFungibleTransfer type")
	}

	return &relayer.Message{
		Source:       sourceID,
		Destination:  uint8(evt.Destination),
		DepositNonce: uint64(evt.DepositNonce),
		ResourceId:   evt.ResourceId,
		Payload: []interface{}{
			[]byte(evt.TokenId),
			[]byte(evt.Recipient),
			[]byte(evt.Metadata),
		},
	}, nil
}

func GenericTransferHandler(sourceID uint8, evtI interface{}) (*relayer.Message, error) {
	evt, ok := evtI.(substrate.EventGenericTransfer)
	if !ok {
		return nil, fmt.Errorf("failed to cast EventGenericTransfer type")
	}
	return &relayer.Message{
		Source:       sourceID,
		Destination:  uint8(evt.Destination),
		DepositNonce: uint64(evt.DepositNonce),
		ResourceId:   evt.ResourceId,
		Payload: []interface{}{
			[]byte(evt.Metadata),
		},
	}, nil
}
