package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//region constants
const (
	StaticToken TokenType = 0
	UserToken   TokenType = 1
	Both        TokenType = 2
)

//end region constants

//region struct
type TokenType int

type APIGate struct {
	headerKey string
	keys      []string
	tokenType TokenType
	config    Config
}

type Config struct {
	TokenList      map[string]string
	AuthURL        string
	ErrorMessages  map[int]string
	Context        func(r *http.Request, authenticatedKey string)
	WithPermission bool
}

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//end region struct

//region functions

func InitGate(headerKey string, keys []string, tokenType TokenType, config Config) *APIGate {
	return &APIGate{headerKey: headerKey, keys: keys, tokenType: tokenType, config: config}
}

func (self *APIGate) StartAuth(c *gin.Context) {

}

//endregion functions
