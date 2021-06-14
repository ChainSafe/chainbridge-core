// package chains

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"reflect"
// 	"testing"
// )

// func TestLoadJSONConfig(t *testing.T) {
// 	file, cfg := createTempConfigFile()
// 	defer os.Remove(file.Name())

// 	res, err := GetConfig(".", file.Name())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !reflect.DeepEqual(res, cfg) {
// 		t.Errorf("did not match\ngot: %+v\nexpected: %+v", res.Chains[0], cfg.Chains[0])
// 	}
// }

// func TestValidateConfig(t *testing.T) {
// 	valid := GeneralChainConfig{
// 		Name:     "chain",
// 		Type:     "ethereum",
// 		Id:       1,
// 		Endpoint: "endpoint",
// 		From:     "0x0",
// 	}

// 	missingType := GeneralChainConfig{
// 		Name:     "chain",
// 		Type:     "",
// 		Id:       1,
// 		Endpoint: "endpoint",
// 		From:     "0x0",
// 	}

// 	missingEndpoint := GeneralChainConfig{
// 		Name:     "chain",
// 		Type:     "ethereum",
// 		Id:       1,
// 		Endpoint: "",
// 		From:     "0x0",
// 	}

// 	missingName := GeneralChainConfig{
// 		Name:     "",
// 		Type:     "ethereum",
// 		Id:       1,
// 		Endpoint: "endpoint",
// 		From:     "0x0",
// 	}

// 	missingId := GeneralChainConfig{
// 		Name:     "",
// 		Type:     "ethereum",
// 		Endpoint: "endpoint",
// 		From:     "0x0",
// 	}

// 	rawConfig := RawChainConfig{}
// 	rawConfig.GeneralChainConfig = valid

// 	cfg := Config{
// 		Chains: []RawChainConfig{rawConfig},
// 	}

// 	err := cfg.validate()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rawConfig.GeneralChainConfig = missingType
// 	cfg = Config{
// 		Chains: []RawChainConfig{rawConfig},
// 	}

// 	err = cfg.validate()
// 	if err == nil {
// 		t.Fatal("must require type field")
// 	}

// 	rawConfig.GeneralChainConfig = missingEndpoint
// 	cfg = Config{
// 		Chains: []RawChainConfig{rawConfig},
// 	}

// 	err = cfg.validate()
// 	if err == nil {
// 		t.Fatal("must require endpoint field")
// 	}

// 	rawConfig.GeneralChainConfig = missingName
// 	cfg = Config{
// 		Chains: []RawChainConfig{rawConfig},
// 	}

// 	err = cfg.validate()
// 	if err == nil {
// 		t.Fatal("must require name field")
// 	}

// 	rawConfig.GeneralChainConfig = missingId
// 	cfg = Config{
// 		Chains: []RawChainConfig{rawConfig},
// 	}

// 	err = cfg.validate()
// 	if err == nil {
// 		t.Fatal("must require chain id field")
// 	}
// }

// func createTempConfigFile() (*os.File, *Config) {
// 	testConfig := NewConfig()
// 	generalCfg := GeneralChainConfig{
// 		Name:     "chain",
// 		Type:     "ethereum",
// 		Id:       1,
// 		Endpoint: "endpoint",
// 		From:     "0x0",
// 	}
// 	ethCfg := RawChainConfig{
// 		GeneralChainConfig: generalCfg,
// 		Opts:               map[string]string{"key": "value"},
// 	}
// 	testConfig.Chains = []RawChainConfig{ethCfg}
// 	tmpFile, err := ioutil.TempFile(".", "*.json")
// 	fmt.Println(tmpFile.Name())
// 	if err != nil {
// 		fmt.Println("Cannot create temporary file", "err", err)
// 		os.Exit(1)
// 	}

// 	f := testConfig.ToJSON(tmpFile.Name())
// 	return f, testConfig
// }

// func createGeneralConfig() GeneralChainConfig {
// 	return GeneralChainConfig{
// 		Name:     "chain",
// 		Type:     "ethereum",
// 		Id:       1,
// 		Endpoint: "endpoint",
// 		From:     "0x0",
// 	}
// }
