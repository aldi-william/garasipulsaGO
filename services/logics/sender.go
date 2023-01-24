package logics

import (
	"os"
	"user/domains/models"
	"user/onesender"
	"user/services/repositories"
	"user/utils"
)

type ISenderService interface {
	GetSender(req models.Sender)
}

type SenderService struct {
	userRepository repositories.IUserRepository
	useTransaction bool
}

func InitSenderService(userRepo repositories.IUserRepository) *SenderService {
	if utils.IsNil(userRepo) {
		userRepo = repositories.InitUserRepository(nil, nil)
	}

	service := SenderService{
		userRepository: userRepo,
		useTransaction: false,
	}
	return &service
}

func (service *SenderService) GetSender(req models.Sender) {
	onesender.ApiUrl = os.Getenv("API_URL_ONESENDER")
	onesender.ApiKey = os.Getenv("API_KEY_ONESENDER")
	onesender.SendTextMessage(req.Sender_Phone, "testing2")
}
