package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		contextValue := context.WithValue(ctx.Request.Context(), "trace_id", traceID)
		ctx.Request = ctx.Request.WithContext(contextValue)
		ctx.Next()

	}
}
