package routes

import (
	"user-management-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	logPath := "../../internal/logs/http.log"
	logRecoveryPath := "../../internal/logs/recovery.log"

	httpLogger := zerolog.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1,    // MB
		MaxBackups: 5,    // number of backup files
		MaxAge:     5,    // days before deletion
		Compress:   true, // disabled by default (compress)
		LocalTime:  true, // use local time in log
	}).With().Timestamp().Logger()

	recoveryLogger := zerolog.New(&lumberjack.Logger{
		Filename:   logRecoveryPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		LocalTime:  true,
	}).With().Timestamp().Logger()

	r.Use(
		middleware.AuthMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.RateLimiterMiddleware(),
	)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}
