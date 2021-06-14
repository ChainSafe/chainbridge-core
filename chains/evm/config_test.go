package evm

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains"
)

func TestParseChainConfig(t *testing.T) {

	generalConfig := createGeneralConfig()

	input := RawEVMConfig{
		GeneralChainConfig: generalConfig,
		Bridge:             "0x1234",
		Erc20Handler:       "0x1234",
		Erc721Handler:      "0x1234",
		GenericHandler:     "0x1234",
		MaxGasPrice:        20,
		GasMultiplier:      1,
		GasLimit:           10,
		Http:               true,
	}

	out, err := parseConfig(&input)
	if err != nil {
		t.Fatal(err)
	}

	expected := EVMConfig{
		GeneralChainConfig: generalConfig,
		Bridge:             "0x1234",
		Erc20Handler:       "0x1234",
		Erc721Handler:      "0x1234",
		GenericHandler:     "0x1234",
		MaxGasPrice:        big.NewInt(20),
		GasMultiplier:      big.NewFloat(1),
		GasLimit:           big.NewInt(10),
		Http:               true,
	}

	if !reflect.DeepEqual(&expected, out) {
		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
	}
}

//TestChainConfigOneContract Tests chain config providing only one contract
func TestChainConfigOneContract(t *testing.T) {

	generalConfig := createGeneralConfig()

	input := RawEVMConfig{
		GeneralChainConfig: generalConfig,
		Bridge:             "0x1234",
		Erc20Handler:       "0x1234",
		MaxGasPrice:        20,
		GasMultiplier:      1,
		GasLimit:           10,
		Http:               true,
	}

	out, err := parseConfig(&input)

	if err != nil {
		t.Fatal(err)
	}

	expected := EVMConfig{
		GeneralChainConfig: generalConfig,
		Bridge:             "0x1234",
		Erc20Handler:       "0x1234",
		MaxGasPrice:        big.NewInt(20),
		GasMultiplier:      big.NewFloat(1),
		GasLimit:           big.NewInt(10),
		Http:               true,
	}

	if !reflect.DeepEqual(&expected, out) {
		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
	}
}

func TestRequiredOpts(t *testing.T) {
	// No opts provided
	input := RawEVMConfig{}

	_, err := parseConfig(&input)

	if err == nil {
		t.Error("config missing chainId field but no error reported")
	}

	// Empty bridgeContract provided
	input = RawEVMConfig{Bridge: ""}

	_, err2 := parseConfig(&input)

	if err2 == nil {
		t.Error("config missing bridge address field but no error reported")
	}

}

// func TestExtraOpts(t *testing.T) {
// 	input := chains.RawChainConfig{
// 		Opts: map[string]string{
// 			"bridge":        "0x1234",
// 			"gasLimit":      "10",
// 			"maxGasPrice":   "20",
// 			"gasMultiplier": "1",
// 			"http":          "true",
// 			"incorrect_opt": "error",
// 		},
// 	}

// 	_, err := parseConfig(&input)

// 	if err == nil {
// 		t.Error("Config should not accept incorrect opts.")
// 	}
// }

func createGeneralConfig() chains.GeneralChainConfig {
	var id uint8 = 1
	return chains.GeneralChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}
}
