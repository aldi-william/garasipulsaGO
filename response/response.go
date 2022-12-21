package response

//region imports
import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type IResponse interface {
	Success(c *gin.Context)
	Error(c *gin.Context, message ...string)
}

func (res Response) Success(c *gin.Context) {
	res.Status = http.StatusOK
	res.Message = http.StatusText(http.StatusOK)
	c.JSON(http.StatusOK, res)
}

func (res Response) Error(c *gin.Context, message ...string) {
	s, m := generateError(message)
	if res.Status == 0 {
		res.Status = s
	}
	res.Message = m
	c.JSON(res.Status, res)
}

// Mapping error messages
func generateError(messages []string) (int, string) {
	var status int
	if messages == nil {
		messages = append(messages, http.StatusText(http.StatusBadRequest))
	}
	switch messages[0] {
	case http.StatusText(http.StatusInternalServerError):
		status = http.StatusInternalServerError
	default:
		status = http.StatusBadRequest

	}
	message := fmt.Sprint(strings.Join(messages, "\n"))
	return status, message

}
