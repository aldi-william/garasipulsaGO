package logics

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"user/constants"
	"user/domains/models"
	"user/services/repositories"
	"user/utils"

	"github.com/sirupsen/logrus"
)

type IUserService interface {
	GetUserByID(id uint) (*models.User, error)
	CreateUser(req models.User) (*models.User, error)
	GetArticles(id string) (*models.Articles, error)
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

func (service *UserService) CreateUser(req models.User) (*models.User, error) {
	var (
		user *models.User
	)

	// fetch user data
	userData, err := service.userRepository.CreateUser(req)
	if err != nil {
		logrus.Error("error [services][logics][user][CreateUser] ", err)
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
