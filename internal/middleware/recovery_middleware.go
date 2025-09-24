package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(recoveryLogger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()

				recoveryLogger.Error().
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Str("panic", fmt.Sprintf("%v", err)).
					Str("stack", string(stack)).
					Msg("panic occoured")

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Please try again later.",
				})
			}
		}()

		ctx.Next()
	}
}
