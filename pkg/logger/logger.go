package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(level string) {
	zerolog.TimeFieldFormat = time.RFC3339

	loglevel, err := zerolog.ParseLevel(level)
	if err != nil {
		loglevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(loglevel)

	log = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Warn() *zerolog.Event {
	return log.Warn()
}
