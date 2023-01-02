package middlewares

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"user/utils"

	"github.com/gin-gonic/gin"
)

func Moota_Signature() gin.HandlerFunc {
	return func(c *gin.Context) {
		Signature := c.GetHeader("signature")
		post_data, _ := ioutil.ReadAll(c.Request.Body)
		secret := os.Getenv("SECRET_MOOTA")
		h := hmac.New(sha1.New, []byte(secret))
		h.Write(post_data)
		s := "sha1="
		s2 := hex.EncodeToString(h.Sum(nil))
		s3 := s + s2
		utils.PrintLogSukses("x-hub-signature", s2)
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
