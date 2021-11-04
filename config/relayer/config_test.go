package relayer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func TestRunValidateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateTestSuite))
}

func (s *ValidateTestSuite) SetupSuite()    {}
func (s *ValidateTestSuite) TearDownSuite() {}
func (s *ValidateTestSuite) SetupTest()     {}
func (s *ValidateTestSuite) TearDownTest()  {}

func (s *ValidateTestSuite) TestValidateInvalidPrometheusPort() {
	relayerConfig := RelayerConfig{
		PrometheusPort: 100000,
	}

	err := relayerConfig.Validate()

	s.NotNil(err)
}

func (s *ValidateTestSuite) TestValidateValidConfig() {
	relayerConfig := RelayerConfig{
		PrometheusPort: 2112,
	}

	err := relayerConfig.Validate()

	s.Nil(err)
}

type GetRelayerConfigTestSuite struct {
	suite.Suite
}

func TestRunGetRelayerConfigTestSuite(t *testing.T) {
	suite.Run(t, new(GetRelayerConfigTestSuite))
}

func (s *GetRelayerConfigTestSuite) SetupSuite()    {}
func (s *GetRelayerConfigTestSuite) TearDownSuite() {}
func (s *GetRelayerConfigTestSuite) SetupTest()     {}
func (s *GetRelayerConfigTestSuite) TearDownTest()  {}

func (s *GetRelayerConfigTestSuite) TestErrorReadingConfig() {
	_, err := GetRelayerConfig("invalid")

	s.NotNil(err)
}

func (s *GetRelayerConfigTestSuite) TestInvalidConfig() {
	data := RelayerConfig{
		PrometheusPort: 0,
	}
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("test.json", file, 0644)

	_, err := GetRelayerConfig("test.json")

	_ = os.Remove("test.json")
	s.NotNil(err)
}

func (s *GetRelayerConfigTestSuite) TestValidConfig() {
	data := RelayerConfig{
		PrometheusPort: 3000,
		PrometheusPath: "/endpoint",
	}
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("test.json", file, 0644)

	config, err := GetRelayerConfig("test.json")

	_ = os.Remove("test.json")
	s.Nil(err)
	s.Equal(config, RelayerConfig{
		PrometheusPort: 3000,
		PrometheusPath: "/endpoint",
	})
}
