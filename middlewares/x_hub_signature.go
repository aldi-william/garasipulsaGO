package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"os"
	"user/utils"

	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func X_HUB_Signature() gin.HandlerFunc {
	return func(c *gin.Context) {
		Signature := c.GetHeader("X-Hub-Signature")
		post_data, _ := ioutil.ReadAll(c.Request.Body)
		secret := os.Getenv("SECRET_TOKEN_MOOTA")
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(post_data)
		s := "sha1="
		s2 := hex.EncodeToString(h.Sum(nil))
		s3 := s + s2
		utils.PrintLogSukses("signature_moota", s2)
		if Signature != s3 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Permission denied",
			})
			return
		}
		c.Next()
	}
}
