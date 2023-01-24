package server

import (
	"user/controllers"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
}

type controllerRoutes struct {
	swaggerControllers     controllers.SwaggerController
	healthControllers      *controllers.HealthController
	userControllers        *controllers.UserController
	transactionControllers *controllers.TransactionController
	paymentControllers     *controllers.PaymentController
	senderControllers      *controllers.SenderController
}

// RegisterRoutes is used to register url routes API
func initControllers() *controllerRoutes {
	return &controllerRoutes{
		swaggerControllers:     controllers.InitSwaggerController(),
		healthControllers:      &controllers.HealthController{},
		userControllers:        controllers.InitUserController(nil),
		transactionControllers: controllers.InitTransactionController(nil),
		paymentControllers:     controllers.InitPaymentController(nil),
		senderControllers:      controllers.InitSenderController(nil),
	}
}

func RegisterRoutes() {
	baseRouter(initControllers())
}
