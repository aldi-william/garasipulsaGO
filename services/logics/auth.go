package logics

import (
	"os"
	"user/domains/models"

	"github.com/dgrijalva/jwt-go"
)

type IAuthService interface {
	GenerateToken(req models.Auth) (string, error)
}

type AuthService struct {
}

func InitAuthService() *AuthService {
	service := AuthService{}
	return &service
}

var SECRET_KEY = []byte(os.Getenv("SECRET_KEY_JWT"))

func (service *AuthService) GenerateToken(req models.Auth) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = req.User_ID
	claim["role"] = req.Role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}
