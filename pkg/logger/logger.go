package logger

import (
	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func NewLogger(config LoggerConfig) *zerolog.Logger {

	logger := zerolog.New().With().Timestamp().Logger()

	return &logger
}
