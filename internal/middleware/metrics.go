package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			promhttp.Handler().ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime).Seconds()
		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
		).Observe(duration)

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			fmt.Sprintf("%d", c.Writer.Status()),
		).Inc()
	}
}
