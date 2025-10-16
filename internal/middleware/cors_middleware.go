package middleware

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		origin := ctx.Request.Header.Get("Origin")
// 		allowedOrigins := []string{
// 			"http://127.0.0.1:5500",
// 			"http://localhost:5500",
// 			"http://127.0.0.1:3000",
// 			"http://localhost:3000",
// 		}

// 		// Check if origin is allowed
// 		for _, allowedOrigin := range allowedOrigins {
// 			if origin == allowedOrigin {
// 				ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 				break
// 			}
// 		}

// 		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Authorization")
// 		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
// 		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")

// 		if ctx.Request.Method == "OPTIONS" {
// 			ctx.AbortWithStatus(http.StatusNoContent)
// 			return
// 		}

// 		ctx.Next()
// 	}
// }
