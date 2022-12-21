package controllers

import (
	"user/domains/entities"
	"user/domains/models"
	"user/response"
	"user/services/logics"
	"user/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	transactionLogic logics.ITransactionService
}

func InitTransactionController(transactionLogic logics.ITransactionService) *TransactionController {
	if utils.IsNil(transactionLogic) {
		transactionLogic = logics.InitTransactionService(nil)
	}
	controller := TransactionController{
		transactionLogic: transactionLogic,
	}

	return &controller
}

func (transactionController *TransactionController) PostTransaction(c *gin.Context) {
	var req models.Transaction
	err := c.BindJSON(&req)
	if err != nil {
		utils.PrintLog("error [controllers][transaction][BindJSON] ", err)
		logrus.Error("error [controllers][transaction][BindJSON] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
		return
	}

	res, err := transactionController.transactionLogic.CreateTransaction(req)
	if err != nil {
		utils.PrintLog("error [controllers][transaction][CreateTransaction] ", err)
		logrus.Error("error [controllers][transaction][CreateTransaction] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}

func (transactionController *TransactionController) PostTransactionPLN(c *gin.Context) {
	var req models.TransactionPLN
	err := c.BindJSON(&req)
	if err != nil {
		utils.PrintLog("error [controllers][transaction][BindJSON] ", err)
		logrus.Error("error [controllers][transaction][BindJSON] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
		return
	}

	res, err := transactionController.transactionLogic.CreateTransactionPLN(req)
	if err != nil {
		utils.PrintLog("error [controllers][transaction][CreateTransaction] ", err)
		logrus.Error("error [controllers][transaction][CreateTransaction] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}

func (transactionController TransactionController) GetWebsocket(ctx *gin.Context) {
	transactionController.transactionLogic.GetWebsocket(ctx)
}

func (transactionController TransactionController) GetTransaction(c *gin.Context) {
	var (
		req []entities.Transactions
	)

	res, err := transactionController.transactionLogic.GetTransaction(req)
	if err != nil {
		utils.PrintLog("error [controllers][transaction][GetTransaction] ", err)
		logrus.Error("error [controllers][transaction][GetTransaction] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}
