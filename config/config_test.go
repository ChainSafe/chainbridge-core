package config_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/config/relayer"
	"github.com/stretchr/testify/suite"
)

type GetConfigTestSuite struct {
	suite.Suite
}

func TestRunGetConfigTestSuite(t *testing.T) {
	suite.Run(t, new(GetConfigTestSuite))
}

func (s *GetConfigTestSuite) SetupSuite()    {}
func (s *GetConfigTestSuite) TearDownSuite() {}
func (s *GetConfigTestSuite) SetupTest()     {}
func (s *GetConfigTestSuite) TearDownTest()  {}

func (s *GetConfigTestSuite) Test_InvalidPath() {
	_, err := config.GetConfig("invalid")

	s.NotNil(err)
}

func (s *GetConfigTestSuite) Test_MissingChainType() {
	data := config.RawConfig{
		ChainConfigs: []map[string]interface{}{{
			"name": "chain1",
		}},
	}
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("test.json", file, 0644)

	_, err := config.GetConfig("test.json")

	_ = os.Remove("test.json")
	s.NotNil(err)
	s.Equal(err.Error(), "Chain 'type' must be provided for every configured chain")
}

func (s *GetConfigTestSuite) Test_InvalidRelayerConfig() {
	data := config.RawConfig{
		RelayerConfig: relayer.RawRelayerConfig{
			LogLevel: "invalid",
		},
		ChainConfigs: []map[string]interface{}{{
			"name": "chain1",
		}},
	}
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("test.json", file, 0644)

	_, err := config.GetConfig("test.json")

	_ = os.Remove("test.json")
	s.NotNil(err)
	s.Equal(err.Error(), "Unknown log level: invalid")
}

func (s *GetConfigTestSuite) Test_ValidConfig() {
	data := config.RawConfig{
		RelayerConfig: relayer.RawRelayerConfig{
			LogLevel: "info",
		},
		ChainConfigs: []map[string]interface{}{{
			"type": "evm",
			"name": "evm1",
		}},
	}
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("test.json", file, 0644)

	actualConfig, err := config.GetConfig("test.json")

	_ = os.Remove("test.json")
	s.Nil(err)
	s.Equal(actualConfig, config.Config{
		RelayerConfig: relayer.RelayerConfig{
			LogLevel:                  1,
			LogFile:                   "",
			OpenTelemetryCollectorURL: "",
		},
		ChainConfigs: []map[string]interface{}{{
			"type": "evm",
			"name": "evm1",
		}},
	})
}
