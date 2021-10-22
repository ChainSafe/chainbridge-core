package calls

import (
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"

	"reflect"
	"testing"
)

func TestToCallArg(t *testing.T) {
	kp, err := secp256k1.GenerateKeypair()
	if err != nil {
		t.Errorf("got an error but didn't expected it")
	}
	address := common.HexToAddress(kp.Address())

	msg := ethereum.CallMsg{
		From:     common.Address{},
		To:       &address,
		Value:    big.NewInt(1),
		Gas:      uint64(21000),
		GasPrice: big.NewInt(3000),
	}
	got := ToCallArg(msg)
	want := map[string]interface{}{
		"from":     msg.From,
		"to":       msg.To,
		"value":    (*hexutil.Big)(msg.Value),
		"gas":      hexutil.Uint64(msg.Gas),
		"gasPrice": (*hexutil.Big)(msg.GasPrice),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v, given %v", got, want, msg)
	}
}
