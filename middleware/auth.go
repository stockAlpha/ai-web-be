package middleware

import (
	"github.com/gin-gonic/gin"
)

func ValidUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		return
	}
}
