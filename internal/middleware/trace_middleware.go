package middleware

import "github.com/gin-gonic/gin"

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next()

	}
}
