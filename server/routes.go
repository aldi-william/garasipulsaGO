package server

import (
	"user/middlewares"

	"github.com/gin-gonic/gin"
)

const (
	ParentRoute      string = "user"
	TransactionRoute string = "transactions"
	PaymentRoute     string = "payment"
	SenderRoute      string = "sender"
	// fill in as the service name
)

func baseRouter(c *controllerRoutes) {
	// m := melody.New()
	// m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// router.GET("/ws", func(c *gin.Context) {
	// 	m.HandleRequest(c.Writer, c.Request)
	// })
	// m.HandleMessage(func(s *melody.Session, msg []byte) {
	// 	m.Broadcast(msg)
	// })
	auth := router.Group(ParentRoute).Use(middlewares.ErrorHandler)
	{
		auth.GET("health-check", c.healthControllers.Status)
		// auth.POST("/", c.userControllers.CreateUser)
		// auth.POST("/login", c.userControllers.Login)
		// // auth.POST("/register", c.userControllers.CreateUser)
		// auth.GET("/get-articles/:id", c.userControllers.GetArticle)
	}

	send := router.Group(SenderRoute)
	{
		send.POST("/callbackonesender", c.senderControllers.GetSender)
	}

	// transaction := router.Group(TransactionRoute).Use(middlewares.ErrorHandler)
	// {
	// 	transaction.POST("/postTransaction", middlewares.RateIPLimiter(), c.transactionControllers.PostTransaction)
	// 	// transaction.GET("/ws", c.transactionControllers.GetWebsocket)
	// 	transaction.GET("/getTransaction", c.transactionControllers.GetTransaction)
	// 	transaction.GET("/getTransactionByGagal", c.transactionControllers.GetTransactionByGagal)
	// 	transaction.POST("/putTransaction", middlewares.AuthMiddleware(), c.transactionControllers.PutTransactionByInvoice)
	// 	transaction.GET("/getTransactionByID/:id", c.transactionControllers.GetTransactionByID)
	// }
	// payment := router.Group(PaymentRoute)
	// {
	// 	whitelistfromMoota := make(map[string]bool)
	// 	whitelistfromDigiflazz := make(map[string]bool)
	// 	whitelistfromMoota["128.199.173.138"] = false
	// 	whitelistfromMoota["103.236.201.178"] = true
	// 	whitelistfromMoota["103.15.226.52"] = true
	// 	whitelistfromMoota["202.80.218.200"] = false
	// 	whitelistfromMoota["::1"] = false
	// 	whitelistfromDigiflazz["52.74.250.133"] = true
	// 	whitelistfromDigiflazz["202.80.219.52"] = false
	// 	whitelistfromDigiflazz["::1"] = false

	// 	payment.POST("/callback", middlewares.IPWhiteList(whitelistfromMoota), middlewares.Moota_Signature(), c.paymentControllers.CallBackFromMoota)
	// 	payment.POST("/callbackfromdigiflazz", middlewares.IPWhiteList(whitelistfromDigiflazz), middlewares.X_HUB_Signature(), c.paymentControllers.CallBackFromDigiFlazz)
	// }

	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", c.swaggerControllers.Swagger)
	}
}
