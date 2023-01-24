package controllers

import (
	"user/domains/models"
	"user/response"
	"user/services/logics"
	"user/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SenderController struct {
	senderLogic logics.ISenderService
}

func InitSenderController(senderLogic logics.ISenderService) *SenderController {
	if utils.IsNil(senderLogic) {
		senderLogic = logics.InitSenderService(nil)
	}
	controller := SenderController{
		senderLogic: senderLogic,
	}

	return &controller
}

func (h SenderController) GetSender(c *gin.Context) {
	var req models.Sender
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logrus.Error("error [controllers][sender][GetSender][ShouldBindJSON] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
		return
	}
	h.senderLogic.GetSender(req)
}
