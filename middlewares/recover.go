package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"user/response"
	"user/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	fmt.Println(c.Errors)
	defer func(c *gin.Context) {
		// panic recovery
		if rec := recover(); rec != nil {
			response := response.Response{Status: http.StatusInternalServerError}
			response.Error(c, http.StatusText(http.StatusInternalServerError))
			utils.PrintLog("Middleware ErrorHandler", fmt.Sprintf("Panic Error, detail : %v \n stacktrace: \n %v", rec, string(debug.Stack())))
		}
	}(c)
	c.Next()
}
