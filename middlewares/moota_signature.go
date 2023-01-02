package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"user/utils"

	"github.com/gin-gonic/gin"
)

func Moota_Signature() gin.HandlerFunc {
	return func(c *gin.Context) {
		Signature := c.GetHeader("signature")
		post_data := fmt.Sprintf("%s", c.Request.Body)
		secret := os.Getenv("SECRET_TOKEN_MOOTA")
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(post_data))
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
