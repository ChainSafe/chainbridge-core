package cli

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func TestWriteCliIODataToFile(t *testing.T) {
	expectedLog := time.Now().Format("02-01|15:00.000 ") + "url: ws://localhost:8545 gasLimit: 6721975 gasPrice: 20000000000 networkid: 0 =>" +
		"\nCalled cli commands: url, gasLimit, gasPrice, networkid,  with values: ws://localhost:8545, 6721975, 20000000000, 0, "

	var EvmRootCLI = &cobra.Command{
		Use:   "test-cli",
		Short: "Test cli",
		Long:  "Testing writing cli io data to file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return EvmRootCLI.Execute()
		},
	}

	BindEVMCLIFlags(EvmRootCLI)
	writeCliIODataToFile(EvmRootCLI)

	data, _ := os.ReadFile(cliLogsFilename)
	logFromFile := string(data)
	if !reflect.DeepEqual(logFromFile, expectedLog) {
		t.Errorf("Logs did not match\ngot: %#v\nexpected: %#v", logFromFile, expectedLog)
	}

	err := os.Remove(cliLogsFilename)
	if err != nil {
		log.Fatal().Err(err)
	}
}
