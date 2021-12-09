package logger

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureLogger configures logger level and assigns an array
// of writers for logger to write to
func ConfigureLogger(l zerolog.Level, writers ...io.Writer) {
	zerolog.SetGlobalLevel(l)
	mw := io.MultiWriter(writers...)
	log.Logger = zerolog.New(mw).With().Timestamp().Logger()
}
