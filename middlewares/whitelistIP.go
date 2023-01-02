package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IPWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !whitelist[c.ClientIP()] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Your IP Address is not allowed to connect our restfull api",
			})
			return
		}
		c.Next()
	}
}
