package logger_test

import (
	"os"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/rs/zerolog/log"
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
	expectedLog := time.Now().Format("02-01|15:00:00.000 ") +
		"Called evm-cli with args: --gasLimit=\"7000000\" --gasPrice=\"25000000000\" --help=\"false\" --jsonWallet=\"test-wallet\" --jsonWalletPassword=\"test-wallet-password\" --networkid=\"0\" --url=\"test-url\" =>\n"

	rootCmdArgs := []string{
		"--url", "test-url",
		"--gasLimit", "7000000",
		"--gasPrice", "25000000000",
		"--networkid", "0x0",
		"--privateKey", "test-private-key",
		"--jsonWallet", "test-wallet",
		"--jsonWalletPassword", "test-wallet-password",
	}

	cli.EvmRootCLI.SetArgs(rootCmdArgs)
	_ = cli.EvmRootCLI.Execute()

	data, _ := os.ReadFile(logger.CliLogsFilename)
	s.Equal(expectedLog, string(data))

	err := os.Remove(logger.CliLogsFilename)
	if err != nil {
		log.Fatal().Err(err)
	}
}
