package substrate

import (
	"reflect"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains"
)

func TestParseChainConfig(t *testing.T) {

	generalConfig := chains.GeneralChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       1,
		Endpoint: "endpoint",
		From:     "0x0",
	}

	input := chains.RawChainConfig{
		GeneralChainConfig: generalConfig,
		Opts: map[string]string{
			"startBlock":      "9999",
			"useExtendedCall": "true",
		},
	}

	out := ParseConfig(&input)

	expected := SubstrateConfig{
		GeneralChainConfig: generalConfig,
		StartBlock:         9999,
		UseExtendedCall:    true,
	}

	if !reflect.DeepEqual(&expected, out) {
		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
	}
}
