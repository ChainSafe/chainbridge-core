package relayer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

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

func (s *GetRelayerConfigTestSuite) TestValidConfig() {
	data := RelayerConfig{
		OpenTelemetryCollectorURL: "http://url.com",
	}
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("test.json", file, 0644)

	config, err := GetRelayerConfig("test.json")

	_ = os.Remove("test.json")
	s.Nil(err)
	s.Equal(config, RelayerConfig{
		OpenTelemetryCollectorURL: "http://url.com",
	})
}
