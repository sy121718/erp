package middleware

import (
	"time"

	"erp-server/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if query != "" {
			path = path + "?" + query
		}

		// 根据状态码使用不同的日志级别
		if status >= 500 {
			log.Error("Server Error",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.Duration("latency", latency),
				zap.String("ip", clientIP),
			)
		} else if status >= 400 {
			log.Warn("Client Error",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.Duration("latency", latency),
				zap.String("ip", clientIP),
			)
		} else {
			log.Info("Request Success",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.Duration("latency", latency),
				zap.String("ip", clientIP),
			)
		}
	}
}
