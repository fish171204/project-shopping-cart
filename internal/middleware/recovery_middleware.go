package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(httpLogger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Please try again later.",
				})
			}
		}()

		ctx.Next()
	}
}
