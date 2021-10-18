package writer

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	CliLogsFilename = "cli_output_data.log"
)

func WriteCliDataToFile(cmd *cobra.Command) {
	file, err := os.OpenFile(CliLogsFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to create cli log file: %v", err))
	}

	defer file.Close()

	currentTimestamp := time.Now().Format("02-01|15:00.000 ")

	var cmdFlagsWithArgs string
	var flags string
	var args string
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name != "privateKey" {
			flags += fmt.Sprintf(`%s, `, flag.Name)
			args += fmt.Sprintf(`%s, `, flag.Value)
			cmdFlagsWithArgs += fmt.Sprintf("--%s=%q ", flag.Name, flag.Value)
		}
	})

	_, err = file.WriteString(
		currentTimestamp +
			fmt.Sprintf("Passed flags: %s with args: %s=>\n", flags, args) +
			fmt.Sprintf("Called %s with args: %s\n", cmd.Name(), cmdFlagsWithArgs))
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to create cli log file: %v", err))
	}
}
