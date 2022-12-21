package controllers

import (
	"user/domains/models"
	"user/response"
	"user/services/logics"
	"user/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userLogic logics.IUserService
}

func InitUserController(userLogic logics.IUserService) *UserController {
	if utils.IsNil(userLogic) {
		userLogic = logics.InitUserService(nil)
	}
	controller := UserController{
		userLogic: userLogic,
	}

	return &controller
}

func (h UserController) GetUser(c *gin.Context) {
	var (
		request = utils.MapRequest(c, &models.BaseRequest{}, []string{"id"})
	)

	res, err := h.userLogic.GetUserByID(request.ID)
	if err != nil {
		logrus.Error("error [controllers][user][GetUserByID] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}

// CreateUser
// @Summary     Create User
// @Description Create User
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       body body     models.User true "Request Body"
// @Success     200  {object} response.Response{data=models.ResponseOnboarding}
// @Failure     400
// @Router      / [post]
func (h UserController) CreateUser(c *gin.Context) {
	var req models.User
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Error("error [controllers][user][BindJSON] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
		return
	}

	res, err := h.userLogic.CreateUser(req)
	if err != nil {
		logrus.Error("error [controllers][user][CreateUser] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}

// GetArticle
// @Summary     Get Article
// @Description Get Article
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       body body     models.User true "Request Body"
// @Success     200  {object} response.Response{data=models.ResponseOnboarding}
// @Failure     400
// @Router      /get-articles [get]
func (h UserController) GetArticle(c *gin.Context) {

	id := c.Param("id")

	res, err := h.userLogic.GetArticles(id)
	if err != nil {
		logrus.Error("error [controllers][user][CreateUser] ", err)
		response := response.Response{}
		response.Error(c, err.Error())
	} else {
		response := response.Response{Data: res}
		response.Success(c)
	}
}
