package server

import (
	"user/middlewares"

	"github.com/gin-gonic/gin"
)

const (
	ParentRoute      string = "user"
	TransactionRoute string = "transactions" // fill in as the service name
)

func baseRouter(c *controllerRoutes) {
	auth := router.Group(ParentRoute).Use(middlewares.ErrorHandler)
	{
		auth.GET("health-check", c.healthControllers.Status)
		auth.POST("/", c.userControllers.CreateUser)
		auth.GET("/get-articles/:id", c.userControllers.GetArticle)
	}

	transaction := router.Group(TransactionRoute)
	{
		transaction.POST("/postTransaction", middlewares.RateIPLimiter(), c.transactionControllers.PostTransaction)
		transaction.POST("/postTransactionPLN", middlewares.RateIPLimiter(), c.transactionControllers.PostTransactionPLN)
		transaction.GET("/ws", c.transactionControllers.GetWebsocket)
		transaction.GET("/getTransaction", c.transactionControllers.GetTransaction)
	}

	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", c.swaggerControllers.Swagger)
	}
}