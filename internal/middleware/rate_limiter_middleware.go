package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lassSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]Client)
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exists := clients[ip]
	// IP does not exist → create new
	if !exists {
		requestSecStr := utils.GetEnv("RATE_LIMITER_REQUEST_SEC", "5")
		brustStr := utils.GetEnv("RATE_LIMITER_REQUEST_BURST", "10")

		requestSec, err := strconv.Atoi(requestSecStr)
		if err != nil {
			panic("invalid RATE_LIMITER_REQUEST_SEC: " + err.Error())
		}

		brust, err := strconv.Atoi(brustStr)
		if err != nil {
			panic("invalid RATE_LIMITER_REQUEST_BURST: " + err.Error())
		}

		limiter := rate.NewLimiter(rate.Limit(requestSec), brust) // 5 request/s , brust : 10 (max), ban đầu 10, hết 10 cấp phát thêm 5 rq mỗi giây
		newClient := &Client{limiter, time.Now()}
		clients[ip] = *newClient

		return limiter
	}

	client.lassSeen = time.Now()
	clients[ip] = client
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lassSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// hey -n 20 -c 1 -H "X-API-Key:(trong .env)" http://localhost:8080/api/v1/users
func RateLimiterMiddleware(rateLimiterLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)

		limiter := getRateLimiter(ip)
		if !limiter.Allow() {

			rateLimiterLogger.Warn().
				Str("method", ctx.Request.Method).
				Str("path", ctx.Request.URL.Path).
				Str("query", ctx.Request.URL.RawQuery).
				Str("client_ip", ctx.ClientIP()).
				Str("user_agent", ctx.Request.UserAgent()).
				Str("referer", ctx.Request.Referer()).
				Str("protocol", ctx.Request.Proto).
				Str("host", ctx.Request.Host).
				Str("remote_addr", ctx.Request.RemoteAddr).
				Str("request_uri", ctx.Request.RequestURI).
				Interface("headers", ctx.Request.Header).
				Msg("rate limiter exceeded")

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many request",
				"message": "Bạn đã gửi quá nhiều request. Hãy thử lại sau",
			})
			return
		}

		ctx.Next()
	}
}
