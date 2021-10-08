package listener

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

func TestErc20EventHandler(t *testing.T) {
	// 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	// construct ERC20 deposit data
	// follows behavior of solidity tests
	// https://github.com/ChainSafe/chainbridge-solidity/blob/develop/test/contractBridge/depositERC20.js#L46-L50
	var calldata []byte
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(2), 32)...)
	calldata = append(calldata, math.PaddedBigBytes(big.NewInt(int64(len(recipientByteSlice))), 32)...)
	calldata = append(calldata, recipientByteSlice...)

	depositLog := &DepositLogs{
		DestinationID:   0,
		ResourceID:      [32]byte{0},
		DepositNonce:    1,
		SenderAddress:   common.HexToAddress("0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"),
		Calldata:        calldata,
		HandlerResponse: []byte{}, // empty []byte
	}

	sourceID := uint8(1)

	message, err := Erc20EventHandler(sourceID, depositLog.DestinationID, depositLog.DepositNonce, depositLog.ResourceID, depositLog.Calldata, depositLog.HandlerResponse)
	if err != nil {
		t.Fatalf("could not generate event handler message: %v", err)
	}

	// TODO: refactor
	amountParsed := calldata[:32]
	// ignore recipientAddressLenParsed: calldata[33:64]
	recipientAddressParsed := calldata[65:]

	expected := &relayer.Message{
		Source:       uint8(1),
		Destination:  depositLog.DestinationID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         relayer.FungibleTransfer,
		Payload: []interface{}{
			amountParsed,
			recipientAddressParsed,
		},
	}

	if !reflect.DeepEqual(message, expected) {
		t.Fatal("ERC20 event handler message does not equal expected message")
	}
}
