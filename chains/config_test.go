package chains

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLoadJSONConfig(t *testing.T) {
	file, cfg := createTempConfigFile()
	defer os.Remove(file.Name())

	res, err := GetConfig(".", file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, cfg) {
		t.Errorf("did not match\ngot: %+v\nexpected: %+v", res.Chains[0], cfg.Chains[0])
	}
}

func TestValidateConfig(t *testing.T) {
	var id uint8 = 1
	valid := GeneralChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}

	missingType := GeneralChainConfig{
		Name:     "chain",
		Type:     "",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}

	missingEndpoint := GeneralChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       &id,
		Endpoint: "",
		From:     "0x0",
	}

	missingName := GeneralChainConfig{
		Name:     "",
		Type:     "ethereum",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}

	missingId := GeneralChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Endpoint: "endpoint",
		From:     "0x0",
	}

	err := valid.Validate()
	if err != nil {
		t.Fatal(err)
	}

	err = missingType.Validate()
	if err == nil {
		t.Fatal("must require type field")
	}

	err = missingEndpoint.Validate()
	if err == nil {
		t.Fatal("must require endpoint field")
	}

	err = missingName.Validate()
	if err == nil {
		t.Fatal("must require name field")
	}

	err = missingId.Validate()
	if err == nil {
		t.Fatal("must require chain id field")
	}

}

func createTempConfigFile() (*os.File, *Config) {
	testConfig := NewConfig()
	var id uint8 = 1
	generalCfg := GeneralChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       &id,
		Endpoint: "endpoint",
		From:     "0x0",
	}
	ethCfg := RawChainConfig{
		GeneralChainConfig: generalCfg,
		Opts:               map[string]string{"key": "value"},
	}
	testConfig.Chains = []RawChainConfig{ethCfg}
	tmpFile, err := ioutil.TempFile(".", "*.json")
	fmt.Println(tmpFile.Name())
	if err != nil {
		fmt.Println("Cannot create temporary file", "err", err)
		os.Exit(1)
	}

	f := testConfig.ToJSON(tmpFile.Name())
	return f, testConfig
}
