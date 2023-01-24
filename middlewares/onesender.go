package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func OnesenderHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		onesenderkey := c.GetHeader("onesender-key")
		secretOnesender := os.Getenv("SECRET_ONE_SENDER_KEY")
		if onesenderkey != secretOnesender {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Permission denied",
			})
			return
		}
		c.Next()
	}
}
