package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"user/controllers"
	"user/domains/models"

	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := models.Users{Email: "test@gmail.com", Password: "test"}
	c.JSON(http.StatusOK, req)
	h := controllers.UserController{}
	h.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, but got %v", w.Code)
	}

}
