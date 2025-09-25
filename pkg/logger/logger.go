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
		Filename:   "../../internal/logs/http.log",
		MaxSize:    1,    // MB
		MaxBackups: 5,    // number of backup files
		MaxAge:     5,    // days before deletion
		Compress:   true, // disabled by default (compress)
		LocalTime:  true, // use local time in log
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}
