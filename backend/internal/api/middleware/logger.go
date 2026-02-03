package middleware

import (
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

// Logger logs the request details and execution time.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[API] %v | %3d | %13v | %15s | %-7s %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			c.ClientIP(),
			method,
			path,
		)
	}
}
