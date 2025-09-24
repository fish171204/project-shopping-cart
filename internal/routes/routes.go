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

	logger := zerolog.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1,    // MB
		MaxBackups: 5,    // number of backup files
		MaxAge:     5,    // days before deletion
		Compress:   true, // disabled by default (compress)
		LocalTime:  true, // use local time in log
	}).With().Timestamp().Logger()

	r.Use(
		middleware.AuthMiddleware(),
		middleware.LoggerMiddleware(logger),
		middleware.RecoveryMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.RateLimiterMiddleware(),
	)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}
