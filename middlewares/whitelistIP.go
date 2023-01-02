package middlewares

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IPWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ip string
		forwardHeader := c.Request.Header.Get("x-forwarded-for")
		firstAddress := strings.Split(forwardHeader, ",")[0]
		if net.ParseIP(strings.TrimSpace(firstAddress)) != nil {
			ip = firstAddress
		} else {
			ip = c.ClientIP()
		}

		if !whitelist[ip] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": ip,
			})
			return
		}
	}
}
