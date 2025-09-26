package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	// Register metrics only once
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// MetricsMiddleware collects request count and latency
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath() // uses route pattern (e.g. /user/:id)

		if path == "" {
			path = c.Request.URL.Path // fallback to raw path
		}

		// Increase counters
		httpRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
