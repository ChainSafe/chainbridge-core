package writer_test

import (
	"os"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/writer"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(WriterTestSuite))
}

func (s *WriterTestSuite) SetupSuite()    {}
func (s *WriterTestSuite) TearDownSuite() {}

func (s *WriterTestSuite) TearDownTest() {}

func (s *WriterTestSuite) TestWriteCliDataToFile() {
	expectedLog := time.Now().Format("02-01|15:00.000 ") +
		"Passed flags: gasLimit, gasPrice, help, jsonWallet, jsonWalletPassword, networkid, url,  with args: 7000000, 25000000000, false, test-wallet, test-wallet-password, 0, test-url, =>\n" +
		"Called evm-cli with args: --gasLimit=\"7000000\" --gasPrice=\"25000000000\" --help=\"false\" --jsonWallet=\"test-wallet\" --jsonWalletPassword=\"test-wallet-password\" --networkid=\"0\" --url=\"test-url\" \n"

	var EvmRootCLI = &cobra.Command{
		Use:   "evm-cli",
		Short: "EVM CLI",
		Long:  "Root command for starting EVM CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			writer.WriteCliDataToFile(cmd)
			return nil
		},
	}

	rootCmdArgs := []string{
		"evm-cli",
		"--url", "test-url",
		"--gasLimit", "7000000",
		"--gasPrice", "25000000000",
		"--networkid", "0x0",
		"--privateKey", "test-private-key",
		"--jsonWallet", "test-wallet",
		"--jsonWalletPassword", "test-wallet-password"}

	EvmRootCLI.SetArgs(rootCmdArgs)
	cli.BindEVMCLIFlags(EvmRootCLI)
	_ = EvmRootCLI.Execute()

	data, _ := os.ReadFile(writer.CliLogsFilename)
	logFromFile := string(data)
	s.Equal(logFromFile, expectedLog)

	err := os.Remove(writer.CliLogsFilename)
	if err != nil {
		log.Fatal().Err(err)
	}
}
