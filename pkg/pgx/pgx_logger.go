package pgx

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type PgxZerologTracer struct {
	Logger         zerolog.Logger
	SlowQueryLimit time.Duration
}

func (t *PgxZerologTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	log.Printf("%+v", data)

	sql, _ := data["sql"].(string)
	args, _ := data["args"].([]any)
	duration, _ := data["time"].(time.Duration)

	baseLogger := t.Logger.With()

	if msg == "Query" && duration > t.SlowQueryLimit {
		logger.Warn().Str("event", "Query").Msg("Slow SQL Query")
		return
	}

}
