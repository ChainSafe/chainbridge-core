package logger_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/cli"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
	EvmRootCLI *cobra.Command
}

func TestLoggerWriteToFile(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) SetupSuite() {
}
func (s *LoggerTestSuite) TearDownSuite() {}

func (s *LoggerTestSuite) TearDownTest() {}

func (s *LoggerTestSuite) TestWriteCliDataToFile() {
	expectedLog := "Called evm-cli with args: --gas-limit=\"7000000\" --gas-price=\"25000000000\" --help=\"false\" --json-wallet=\"test-wallet\" --json-wallet-password=\"test-wallet-password\" --network=\"0\" --prepare=\"false\" --private-key=\"test-private-key\" --url=\"test-url\" =>\n"

	rootCmdArgs := []string{
		"--url", "test-url",
		"--gas-limit", "7000000",
		"--gas-price", "25000000000",
		"--network", "0x0",
		"--private-key", "test-private-key",
		"--json-wallet", "test-wallet",
		"--json-wallet-password", "test-wallet-password",
	}

	cli.EvmRootCLI.SetArgs(rootCmdArgs)
	_ = cli.EvmRootCLI.Execute()

	data, _ := os.ReadFile(logger.CliLogsFilename)
	logParts := strings.SplitN(string(data), " ", 2)
	s.Equal(expectedLog, logParts[1])
	s.True(regexp.Match("[0-9]{2}-[0-9]{2}|[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3}", []byte(logParts[0])))

	err := os.Remove(logger.CliLogsFilename)
	s.Nil(err)
}
