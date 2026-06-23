package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"gomall-lite-api/internal/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		attrs := []any{
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		log := logger.Default()
		if status >= 500 {
			log.Error("request completed", attrs...)
			return
		}
		if status >= 400 {
			log.Warn("request completed", attrs...)
			return
		}
		log.Log(c.Request.Context(), slog.LevelInfo, "request completed", attrs...)
	}
}
