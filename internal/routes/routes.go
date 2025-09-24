package routes

import (
	"user-management-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	r.Use(
		middleware.AuthMiddleware(),
		middleware.LoggerMiddleware(),
		middleware.RecoveryMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.RateLimiterMiddleware(),
	)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}
