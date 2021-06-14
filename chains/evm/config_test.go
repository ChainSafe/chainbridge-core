// package evm

// import (
// 	"math/big"
// 	"reflect"
// 	"testing"

// 	"github.com/ChainSafe/chainbridge-core/chains"
// )

// func TestParseChainConfig(t *testing.T) {

// 	generalConfig := createGeneralConfig()

// 	input := chains.RawChainConfig{
// 		GeneralChainConfig: generalConfig,
// 		Opts: map[string]string{
// 			"bridge":         "0x1234",
// 			"erc20Handler":   "0x1234",
// 			"erc721Handler":  "0x1234",
// 			"genericHandler": "0x1234",
// 			"gasLimit":       "10",
// 			"gasMultiplier":  "1",
// 			"maxGasPrice":    "20",
// 			"http":           "true",
// 		},
// 	}

// 	out, err := parseConfig(&input)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := EVMConfig{
// 		GeneralChainConfig: generalConfig,
// 		Bridge:             "0x1234",
// 		Erc20Handler:       "0x1234",
// 		Erc721Handler:      "0x1234",
// 		GenericHandler:     "0x1234",
// 		MaxGasPrice:        big.NewInt(20),
// 		GasMultiplier:      big.NewFloat(1),
// 		GasLimit:           big.NewInt(10),
// 		Http:               true,
// 	}

// 	if !reflect.DeepEqual(&expected, out) {
// 		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
// 	}
// }

// //TestChainConfigOneContract Tests chain config providing only one contract
// func TestChainConfigOneContract(t *testing.T) {

// 	generalConfig := createGeneralConfig()

// 	input := chains.RawChainConfig{
// 		GeneralChainConfig: generalConfig,
// 		Opts: map[string]string{
// 			"bridge":        "0x1234",
// 			"erc20Handler":  "0x1234",
// 			"gasLimit":      "10",
// 			"maxGasPrice":   "20",
// 			"gasMultiplier": "1",
// 			"http":          "true",
// 		},
// 	}

// 	out, err := ParseConfig(&input)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := EVMConfig{
// 		GeneralChainConfig: generalConfig,
// 		Bridge:             "0x1234",
// 		Erc20Handler:       "0x1234",
// 		MaxGasPrice:        big.NewInt(20),
// 		GasMultiplier:      big.NewFloat(1),
// 		GasLimit:           big.NewInt(10),
// 		Http:               true,
// 	}

// 	if !reflect.DeepEqual(&expected, out) {
// 		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
// 	}
// }

// func TestRequiredOpts(t *testing.T) {
// 	// No opts provided
// 	input := chains.RawChainConfig{
// 		Opts: map[string]string{},
// 	}

// 	_, err := ParseConfig(&input)

// 	if err == nil {
// 		t.Error("config missing chainId field but no error reported")
// 	}

// 	// Empty bridgeContract provided
// 	input = chains.RawChainConfig{
// 		Opts: map[string]string{"bridge": ""},
// 	}

// 	_, err2 := ParseConfig(&input)

// 	if err2 == nil {
// 		t.Error("config missing chainId field but no error reported")
// 	}

// }

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

// 	_, err := ParseConfig(&input)

// 	if err == nil {
// 		t.Error("Config should not accept incorrect opts.")
// 	}
// }

// func createGeneralConfig() chains.GeneralChainConfig {
// 	return chains.GeneralChainConfig{
// 		Name:     "chain",
// 		Type:     "ethereum",
// 		Id:       1,
// 		Endpoint: "endpoint",
// 		From:     "0x0",
// 	}
// }
