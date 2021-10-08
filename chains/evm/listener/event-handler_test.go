package listener

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// moved here to circumvent import cycle recursion
// https://github.com/ChainSafe/chainbridge-core/blob/main/chains/evm/calls/utils.go#L108-L114
func constructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

func TestErc20EventHandler(t *testing.T) {
	recipientByteSlice := []byte{241, 229, 143, 177, 119, 4, 194, 218, 132, 121, 165, 51, 249, 250, 212, 173, 9, 147, 202, 107}

	// construct ERC20 deposit data
	calldata := constructErc20DepositData(
		recipientByteSlice, // 0xf1e58fb17704c2da8479a533f9fad4ad0993ca6b
		big.NewInt(2),      // 2 tokens
	)

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
	recipientParsed := calldata[33:64]

	expected := &relayer.Message{
		Source:       uint8(1),
		Destination:  depositLog.DestinationID,
		DepositNonce: depositLog.DepositNonce,
		ResourceId:   depositLog.ResourceID,
		Type:         relayer.FungibleTransfer,
		Payload: []interface{}{
			amountParsed,
			recipientParsed,
		},
	}

	if !reflect.DeepEqual(message, expected) {
		t.Fatal("ERC20 event handler message does not equal expected message")
	}
}
