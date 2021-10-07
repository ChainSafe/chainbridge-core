package listener

import (
	"reflect"
	"testing"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
)

func TestErc20EventHandler(t *testing.T) {
	depositLog := &DepositLogs{
		DestinationID:   0,
		ResourceID:      [32]byte{0},
		DepositNonce:    1,
		Address:         common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Calldata:        []byte{123}, // bytes rep. of the amount
		HandlerResponse: []byte{},    // empty []byte
	}

	sourceID := uint8(1)

	message, err := Erc20EventHandler(sourceID, depositLog.DestinationID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Address, depositLog.Calldata, depositLog.HandlerResponse)
	if err != nil {
		t.Fatalf("could not generate event handler message: %v", err)
	}

	expected := &relayer.Message{
		Source:       uint8(1),
		Destination:  depositLog.DestinationID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         relayer.FungibleTransfer,
		Payload: []interface{}{
			depositLog.Calldata,
			depositLog.Address,
			depositLog.HandlerResponse,
		},
	}

	if !reflect.DeepEqual(message, expected) {
		t.Fatal("ERC20 event handler message does not equal expected message")
	}
}
