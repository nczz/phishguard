package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	loginLimiters = sync.Map{}
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// LoginRateLimit allows 5 attempts per minute per IP.
func LoginRateLimit() gin.HandlerFunc {
	// Cleanup stale entries every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			loginLimiters.Range(func(key, value any) bool {
				if time.Since(value.(*ipLimiter).lastSeen) > 5*time.Minute {
					loginLimiters.Delete(key)
				}
				return true
			})
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		v, _ := loginLimiters.LoadOrStore(ip, &ipLimiter{
			limiter: rate.NewLimiter(rate.Every(time.Minute/5), 5), // 5 req/min, burst 5
		})
		entry := v.(*ipLimiter)
		entry.lastSeen = time.Now()

		if !entry.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "登入嘗試過於頻繁，請稍後再試"})
			return
		}
		c.Next()
	}
}
