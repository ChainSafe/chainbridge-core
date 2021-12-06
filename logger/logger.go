package logger

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogger(l zerolog.Level, writers ...io.Writer) {
	zerolog.SetGlobalLevel(l)
	mw := io.MultiWriter(writers...)
	log.Logger = zerolog.New(mw).With().Timestamp().Logger()
}
