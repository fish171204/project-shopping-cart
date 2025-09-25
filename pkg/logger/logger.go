package logger

import (
	"io"
	"os"
	"time"

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
	IsDev      string
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	var writer io.Writer

	if config.IsDev == "development" {
		writer = os.Stdout
	} else {
		writer = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,    // MB
			MaxBackups: config.MaxBackups, // number of backup files
			MaxAge:     config.MaxAge,     // days before deletion
			Compress:   config.Compress,   // disabled by default (compress)
		}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}
