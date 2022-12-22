package controllers

import (
	"user/domains/models"
	"user/services/logics"
	"user/utils"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentLogic logics.IPaymentService
}

func InitPaymentController(paymentLogic logics.IPaymentService) *PaymentController {
	if utils.IsNil(paymentLogic) {
		paymentLogic = logics.InitPaymentService(nil, nil)
	}
	controller := PaymentController{
		paymentLogic: paymentLogic,
	}

	return &controller
}

func (paymentController *PaymentController) CallBackFromMoota(ctx *gin.Context) {
	var req models.MootaCallback
	err := ctx.BindJSON(&req)
	if err != nil {
		utils.PrintLog("[controllers][CallBackFromMoota][BindJSON]", err)
	}

	paymentController.paymentLogic.CallBackFromMoota(req)
}
