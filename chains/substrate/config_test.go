package substrate

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains"
)

func TestParseChainConfig(t *testing.T) {
	var id uint8 = 1
	generalConfig := chains.GeneralChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}

	input := RawSubstrateConfig{
		GeneralChainConfig: generalConfig,
		StartBlock:         9999,
		UseExtendedCall:    true,
	}

	out := parseConfig(&input)

	expected := SubstrateConfig{
		GeneralChainConfig: generalConfig,
		StartBlock:         big.NewInt(9999),
		UseExtendedCall:    true,
	}

	if !reflect.DeepEqual(&expected, out) {
		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
	}
}
