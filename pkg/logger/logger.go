package logger

import (
	"io"

	"github.com/natefinch/lumberjack"
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
	var writer io.Writer

	writer = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,    // MB
		MaxBackups: config.MaxBackups, // number of backup files
		MaxAge:     config.MaxAge,     // days before deletion
		Compress:   config.Compress,   // disabled by default (compress)
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}
