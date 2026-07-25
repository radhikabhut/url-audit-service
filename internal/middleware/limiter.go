package middleware

import (
	"net/http"
	"url-audit-service/internal/limiter"
	"url-audit-service/internal/response"

	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		allowed, err := l.Allow(c.Request.Context(), ip)
		if err != nil || !allowed {
			response.Error(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
