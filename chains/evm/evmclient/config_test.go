package evmclient

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/ChainSafe/chainbridge-core/config/chain"
)

func TestLoadJSONConfig(t *testing.T) {
	file, cfg := createTempConfigFile()
	defer os.Remove(file.Name())

	res, err := GetConfig(".", file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, cfg) {
		t.Errorf("did not match\ngot: %+v\nexpected: %+v", res, cfg)
	}
}

func TestParseChainConfig(t *testing.T) {

	generalConfig := createGeneralConfig()

	input := RawEVMConfig{
		RawSharedEVMConfig: chain.RawSharedEVMConfig{
			GeneralChainConfig: generalConfig,
			Bridge:             "0x1234",
			Erc20Handler:       "0x1234",
			Erc721Handler:      "0x1234",
			GenericHandler:     "0x1234",
			MaxGasPrice:        20,
			GasMultiplier:      1,
			GasLimit:           10,
			StartBlock:         9999,
			BlockConfirmations: 10,
		},
	}

	out, err := ParseConfig(&input)
	if err != nil {
		t.Fatal(err)
	}

	expected := EVMConfig{
		SharedEVMConfig: chain.SharedEVMConfig{
			GeneralChainConfig: generalConfig,
			Bridge:             "0x1234",
			Erc20Handler:       "0x1234",
			Erc721Handler:      "0x1234",
			GenericHandler:     "0x1234",
			MaxGasPrice:        big.NewInt(20),
			GasMultiplier:      big.NewFloat(1),
			GasLimit:           big.NewInt(10),
			StartBlock:         big.NewInt(9999),
			BlockConfirmations: big.NewInt(10),
		},
	}

	if !reflect.DeepEqual(&expected, out) {
		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
	}
}

// TestParseChainConfigWithNoBlockConfirmations Tests chain config without block confirmations
func TestParseChainConfigWithNoBlockConfirmations(t *testing.T) {
	generalConfig := createGeneralConfig()

	input := RawEVMConfig{
		RawSharedEVMConfig: chain.RawSharedEVMConfig{
			GeneralChainConfig: generalConfig,
			Bridge:             "0x1234",
			Erc20Handler:       "0x1234",
			Erc721Handler:      "0x1234",
			GenericHandler:     "0x1234",
			MaxGasPrice:        20,
			GasMultiplier:      1,
			GasLimit:           10,
			StartBlock:         9999,
		},
	}

	out, err := ParseConfig(&input)

	if err != nil {
		t.Fatal(err)
	}

	expected := EVMConfig{
		SharedEVMConfig: chain.SharedEVMConfig{
			GeneralChainConfig: generalConfig,
			Bridge:             "0x1234",
			Erc20Handler:       "0x1234",
			Erc721Handler:      "0x1234",
			GenericHandler:     "0x1234",
			MaxGasPrice:        big.NewInt(20),
			GasMultiplier:      big.NewFloat(1),
			GasLimit:           big.NewInt(10),
			StartBlock:         big.NewInt(9999),
			BlockConfirmations: big.NewInt(10),
		},
	}

	if !reflect.DeepEqual(&expected, out) {
		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
	}
}

//TestChainConfigOneContract Tests chain config providing only one contract
func TestChainConfigOneContract(t *testing.T) {

	generalConfig := createGeneralConfig()

	input := RawEVMConfig{
		RawSharedEVMConfig: chain.RawSharedEVMConfig{
			GeneralChainConfig: generalConfig,
			Bridge:             "0x1234",
			Erc20Handler:       "0x1234",
			MaxGasPrice:        20,
			GasMultiplier:      1,
			GasLimit:           10,
		},
	}

	out, err := ParseConfig(&input)

	if err != nil {
		t.Fatal(err)
	}

	expected := EVMConfig{
		SharedEVMConfig: chain.SharedEVMConfig{
			GeneralChainConfig: generalConfig,
			Bridge:             "0x1234",
			Erc20Handler:       "0x1234",
			MaxGasPrice:        big.NewInt(20),
			GasMultiplier:      big.NewFloat(1),
			GasLimit:           big.NewInt(10),
			StartBlock:         big.NewInt(0),
			BlockConfirmations: big.NewInt(10),
		},
	}

	if !reflect.DeepEqual(&expected, out) {
		t.Fatalf("Output not expected.\n\tExpected: %#v\n\tGot: %#v\n", &expected, out)
	}
}

func TestRequiredOpts(t *testing.T) {
	// No opts provided
	input := RawEVMConfig{}

	_, err := ParseConfig(&input)

	if err == nil {
		t.Error("config missing domainId field but no error reported")
	}

	// Empty bridgeContract provided
	input = RawEVMConfig{RawSharedEVMConfig: chain.RawSharedEVMConfig{Bridge: ""}}

	_, err = ParseConfig(&input)

	if err == nil {
		t.Error("config missing bridge address field but no error reported")
	}

}

func createGeneralConfig() chain.GeneralChainConfig {
	var id uint8 = 1
	return chain.GeneralChainConfig{
		Name:     "chain",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}
}

func createTempConfigFile() (*os.File, *RawEVMConfig) {
	generalCfg := createGeneralConfig()
	ethCfg := RawEVMConfig{
		RawSharedEVMConfig: chain.RawSharedEVMConfig{
			GeneralChainConfig: generalCfg,
			Bridge:             "0x1234"},
	}
	tmpFile, err := ioutil.TempFile(".", "*.json")
	//fmt.Println(*generalCfg.Id)
	if err != nil {
		fmt.Println("Cannot create temporary file", "err", err)
		os.Exit(1)
	}

	f := ethCfg.ToJSON(tmpFile.Name())
	return f, &ethCfg
}
