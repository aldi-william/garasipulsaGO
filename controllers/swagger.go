package controllers

import (
	"fmt"
	"os"
	"user/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title       User API
// @version     1.0
// @description User API Documentation

// @host     localhost:7000
// @BasePath /user

type SwaggerController struct{}

func InitSwaggerController() SwaggerController {
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
	docs.SwaggerInfo.BasePath = "/user"
	return SwaggerController{}
}

func (SwaggerController) Swagger(c *gin.Context) {
	ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
}
