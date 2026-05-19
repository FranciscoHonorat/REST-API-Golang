package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("10000-S")
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		context, err := instance.Get(c.Request.Context(), ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests", "retry_after": context.Reset})
			c.Abort()
			return
		}

		c.Next()
	}
}
