package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Moota_Signature() gin.HandlerFunc {
	return func(c *gin.Context) {
		Signature := c.GetHeader("signature")
		post_data, _ := io.ReadAll(c.Request.Body)
		secret := os.Getenv("SECRET_TOKEN_MOOTA")
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(post_data)
		c.Request.Body = io.NopCloser(bytes.NewReader(post_data))
		s2 := hex.EncodeToString(h.Sum(nil))
		fmt.Println(s2)
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
