package controllers

import (
	"user/response"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

// Health Status is used to check status of the service
// Status
// @Summary     Status
// @Description Status
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       body body     models.User true "Request Body"
// @Success     200  {object} response.Response{data=models.ResponseOnboarding}
// @Failure     400
// @Router      /health-check [get]
func (h HealthController) Status(c *gin.Context) {
	var resp response.IResponse = response.Response{}
	resp.Success(c)
}
