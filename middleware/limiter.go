package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(redisClient *redis.Client, maxAttempts int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		key := "login_attempts:" + ip

		attempts, _ := redisClient.Get(context.Background(), key).Int()
		if attempts >= maxAttempts {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many login attempts"})
			ctx.Abort()
			return
		}

		redisClient.Incr(context.Background(), key)
		redisClient.Expire(context.Background(), key, time.Minute*5)
		ctx.Next()
	}
}

func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		for _, ip := range allowedIPs {
			if clientIP == ip {
				// IP is allowed, continue processing the request
				c.Next()
				return
			}
		}
		// If the IP is not allowed, return an unauthorized error
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: IP not allowed"})
		c.Abort()
	}
}
