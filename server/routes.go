package server

import (
	"user/middlewares"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

const (
	ParentRoute      string = "user"
	TransactionRoute string = "transactions"
	PaymentRoute     string = "payment"
	// fill in as the service name
)

func baseRouter(c *controllerRoutes) {
	m := melody.New()
	router.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})
	auth := router.Group(ParentRoute).Use(middlewares.ErrorHandler)
	{
		auth.GET("health-check", c.healthControllers.Status)
		auth.POST("/", c.userControllers.CreateUser)
		auth.GET("/get-articles/:id", c.userControllers.GetArticle)
	}

	transaction := router.Group(TransactionRoute).Use(middlewares.ErrorHandler)
	{
		transaction.POST("/postTransaction", middlewares.RateIPLimiter(), c.transactionControllers.PostTransaction)
		// transaction.GET("/ws", c.transactionControllers.GetWebsocket)
		transaction.GET("/getTransaction", c.transactionControllers.GetTransaction)
	}
	payment := router.Group(PaymentRoute)
	{
		whitelistfromMoota := make(map[string]bool)
		whitelistfromDigiflazz := make(map[string]bool)
		whitelistfromMoota["128.199.173.138"] = true
		whitelistfromMoota["202.80.219.52"] = true
		whitelistfromMoota["::1"] = true
		whitelistfromDigiflazz["52.74.250.133"] = true

		payment.POST("/callback", middlewares.IPWhiteList(whitelistfromMoota), c.paymentControllers.CallBackFromMoota)
		payment.POST("/callbackfromdigiflazz", middlewares.IPWhiteList(whitelistfromDigiflazz), c.paymentControllers.CallBackFromDigiFlazz)
	}

	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", c.swaggerControllers.Swagger)
	}
}
