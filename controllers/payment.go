package controllers

import (
	"user/domains/models"
	"user/response"
	"user/services/logics"
	"user/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func (paymentController *PaymentController) CallBackFromMoota(c *gin.Context) {
	var req []models.MootaCallback
	err := c.Bind(&req)
	if err != nil {
		utils.PrintLog("error [controllers][CallBackFromMoota][BindJSON]", err)
		logrus.Error("error [controllers][CallBackFromMoota][BindJSON] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
		return
	}

	res, err := paymentController.paymentLogic.CallBackFromMoota(req)
	if err != nil {
		utils.PrintLog("error [controllers][CallBackFromMoota][paymentLogic] ", err)
		logrus.Error("error [controllers][CallBackFromMoota][paymentLogic] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}

func (paymentController *PaymentController) CallBackFromDigiFlazz(c *gin.Context) {
	var req models.DigiFlazzData
	err := c.BindJSON(&req)
	if err != nil {
		utils.PrintLog("error [controllers][CallBackFromMoota][BindJSON] ", err)
		logrus.Error("error [controllers][CallBackFromMoota][BindJSON] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
		return
	}

	res, err := paymentController.paymentLogic.CallBackFromDigiflazz(req)
	if err != nil {
		utils.PrintLog("error [controllers][CallBackFromMoota][paymentLogic] ", err)
		logrus.Error("error [controllers][CallBackFromMoota][paymentLogic] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}
