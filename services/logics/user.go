package logics

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"user/constants"
	"user/domains/entities"
	"user/domains/models"
	"user/services/repositories"
	"user/utils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUserByID(id uint) (*models.User, error)
	GetArticles(id string) (*models.Articles, error)
	GetLogin(req models.Users) (*models.Users, error)
	RegisterUser(req models.Users) (*entities.Users, error)
}

type UserService struct {
	userRepository repositories.IUserRepository
	useTransaction bool
}

func InitUserService(userRepo repositories.IUserRepository) *UserService {
	if utils.IsNil(userRepo) {
		userRepo = repositories.InitUserRepository(nil, nil)
	}

	service := UserService{
		userRepository: userRepo,
		useTransaction: false,
	}
	return &service
}

func (service *UserService) GetUserByID(id uint) (*models.User, error) {
	var (
		user *models.User
	)

	// fetch user data
	userData, err := service.userRepository.GetUserByID(id)
	if err != nil {
		logrus.Error("error [services][logics][user][GetUserByID] ", err)
		return nil, errors.New(constants.UserNotFoundErr)
	}
	_ = utils.AutoMap(userData, &user)

	return user, nil
}

func (service *UserService) GetArticles(id string) (*models.Articles, error) {

	var (
		baseUrl  = os.Getenv("URL_ARTICLE")
		url      = baseUrl + "/posts/" + id
		articles = models.Articles{}
	)

	resp, err := utils.CallAPI(http.MethodGet, url, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &articles)
	if err != nil {
		return nil, err
	}

	return &articles, nil
}

func (service *UserService) GetLogin(req models.Users) (*models.Users, error) {

	email := req.Email
	password := req.Password

	userByEmail, err := service.userRepository.FindUserByEmail(email)
	if err != nil {
		utils.PrintLog("error [services][logics][GetLogin][FindUserByEmail] ", err)
		logrus.Error("error [services][logics][GetLogin][FindUSerByEmail] ", err)
		return nil, errors.New(constants.UserNotFoundErr)
	}
	res := models.Users{}
	res.ID = userByEmail.ID
	res.Email = userByEmail.Email
	res.Name = userByEmail.Name
	res.Role = userByEmail.Role

	token, err := utils.GenerateToken(res)

	res.Token = token

	if userByEmail.ID == 0 {
		utils.PrintLog("error [services][logics][GetLogin][user not found] ", err)
		logrus.Error("error [services][logics][GetLogin][user not found] ", err)
		return &res, errors.New(constants.UserNotFoundErr)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userByEmail.Password), []byte(password))
	if err != nil {
		utils.PrintLog("error [services][logics][GetLogin][CompareHashAndPassword] ", err)
		logrus.Error("error [services][logics][GetLogin][CompareHashAndPassword] ", err)
		return &res, err
	}

	return &res, nil
}

func (service *UserService) RegisterUser(req models.Users) (*entities.Users, error) {
	user := models.Users{}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, errors.New(constants.FailedToGeneratePassword)
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Password = string(passwordHash)
	user.Role = "administrator"
	newUser, err := service.userRepository.CreateUser(user)
	if err != nil {
		return nil, errors.New(constants.FailedToCreatePassword)
	}
	utils.PrintLogSukses("SUCCESS TO CREATE USER", newUser)
	return nil, nil
}
