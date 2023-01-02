package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"user/utils"

	"github.com/gin-gonic/gin"
)

func Moota_Signature() gin.HandlerFunc {
	return func(c *gin.Context) {
		Signature := c.GetHeader("signature")
		post_data, _ := io.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		secret := os.Getenv("SECRET_TOKEN_MOOTA")
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(post_data)
		s2 := hex.EncodeToString(h.Sum(nil))
		utils.PrintLogSukses("signature from moota", s2)
		if Signature != s2 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Permission denied",
			})
			return
		}
		c.Next()
	}
}
