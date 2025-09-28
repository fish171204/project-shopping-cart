package pgx

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
)

type PgxZerologTracer struct {
}

func (t *PgxZerologTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any)
