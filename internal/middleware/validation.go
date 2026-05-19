package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const MaxPayloadSize = 1048576

func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			contentType := c.Request.Header.Get("Content-Type")
			if contentType == "" || !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}

			if c.Request.ContentLength > MaxPayloadSize {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": fmt.Sprintf("Payload too large. Max size is %d bytes", MaxPayloadSize),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
