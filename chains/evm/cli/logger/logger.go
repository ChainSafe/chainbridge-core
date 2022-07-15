package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/ChainSafe/sygma-core/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

var (
	CliLogsFilename = "cli_output_data.log"
	Now             = time.Now
)

func LoggerMetadata(cmdName string, flagSet *pflag.FlagSet) {

	currentTimestamp := Now().Format("02-01|15:00:00.000 ")

	file, err := os.OpenFile(CliLogsFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to create cli log file: %v", err))
	}

	var cmdFlagsWithArgs string
	flagSet.VisitAll(func(flag *pflag.Flag) {
		if flag.Name != "privateKey" {
			cmdFlagsWithArgs += fmt.Sprintf("--%s=%q ", flag.Name, flag.Value)
		}
	})

	_, err = file.WriteString(
		currentTimestamp +
			fmt.Sprintf("Called %s with args: %s=>\n", cmdName, cmdFlagsWithArgs))

	if err != nil {
		log.Error().Err(fmt.Errorf("failed to write to log file: %v", err))
	}

	// PartsExclude - omit log level and execution time from final log
	logConsoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{"level", "time"}}
	logFileWriter := zerolog.ConsoleWriter{Out: file, PartsExclude: []string{"level", "time"}}
	logger.ConfigureLogger(zerolog.DebugLevel, logConsoleWriter, logFileWriter)
}
