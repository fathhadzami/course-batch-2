package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func WithTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		newCtx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(newCtx)
	}
}
