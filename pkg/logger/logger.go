package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(serviceName string) *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	return &Logger{
		Logger: zerolog.New(output).
			With().
			Str("service", serviceName).
			Logger(),
	}
}
