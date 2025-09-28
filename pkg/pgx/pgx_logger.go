package pgx

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type PgxZerologTracer struct {
	Logger         zerolog.Logger
	SlowQueryLimit time.Duration
}

type QueryInfo struct {
	QueryName     string
	OperationType string
	CleanSQL      string
	OriginalSQL   string
}

var (
	sqlcNameRegex = regexp.MustCompile(`-- name:\s*(\w+)\s*:(\w+)`)
	spaceRegex    = regexp.MustCompile(`\s+`)
	commentRegex  = regexp.MustCompile(`-- [^\r\n]*`)
)

// OriginalSQL: -- name: CreateUser :one .... [0] - Full match
// QueryName: CreateUser					  [1] - Group 1 (\w+)
// OperationType: one						  [2] - Group 2 (\w+)

func parseSQL(sql string) QueryInfo {
	info := QueryInfo{
		OriginalSQL: sql,
	}

	if matches := sqlcNameRegex.FindStringSubmatch(sql); len(matches) == 3 {
		info.QueryName = matches[1]
		info.OperationType = strings.ToUpper(matches[2])
	}

	cleanSQL := commentRegex.ReplaceAllString(sql, "")
	cleanSQL = strings.TrimSpace(cleanSQL)
	cleanSQL = spaceRegex.ReplaceAllString(cleanSQL, " ")
	info.CleanSQL = cleanSQL

	return info
}

func (t *PgxZerologTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {

	sql, _ := data["sql"].(string)
	args, _ := data["args"].([]any)
	duration, _ := data["time"].(time.Duration)

	queryInfo := parseSQL(sql)

	baseLogger := t.Logger.With().
		Dur("duration", duration).
		Str("sql_original", queryInfo.OriginalSQL).
		Str("sql", queryInfo.CleanSQL).
		Str("query_name", queryInfo.QueryName).
		Str("operation", queryInfo.OperationType).
		Interface("args", args)

	logger := baseLogger.Logger()

	if msg == "Query" && duration > t.SlowQueryLimit {
		logger.Warn().Str("event", "Query").Msg("Slow SQL Query")
		return
	}

	if msg == "Query" {
		logger.Info().Str("event", "Query").Msg("Executed SQL")
	}
}
