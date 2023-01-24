package logics

import (
	"net/http"
	"os"
	"user/domains/models"
	"user/onesender"
	"user/services/repositories"
	"user/utils"

	"github.com/sirupsen/logrus"
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
	// onesender.SendTextMessage(req.Sender_Phone, "testing2")
	url := os.Getenv("BASE_URL_API_ONESENDER")
	jsonData := models.APISenderWithButton{}
	jsonData.Recipient_type = "individual"
	jsonData.To = req.Sender_Phone
	jsonData.Type = "interactive"
	jsonData.Interactive.Type = "button"
	jsonData.Interactive.Header.Text = "ini adalah header button"
	jsonData.Interactive.Body.Text = "Test button with header"
	jsonData.Interactive.Footer.Text = "Pilihan Jumlah Donasi"

	// button 1
	button := models.Button{}
	button.Type = "button"
	button.Reply.ID = "rp25000"
	button.Reply.Title = "Rp 25.000,-"
	jsonData.Interactive.Action.Buttons = append(jsonData.Interactive.Action.Buttons, button)
	// button 2
	button2 := models.Button{}
	button2.Type = "button"
	button2.Reply.ID = "rp50000"
	button2.Reply.Title = "Rp 50.000,-"
	jsonData.Interactive.Action.Buttons = append(jsonData.Interactive.Action.Buttons, button2)
	// button 3
	button3 := models.Button{}
	button3.Type = "button"
	button3.Reply.ID = "rp100000"
	button3.Reply.Title = "Rp 100.000,-"
	jsonData.Interactive.Action.Buttons = append(jsonData.Interactive.Action.Buttons, button3)

	headers := map[string]string{
		"Authorization": os.Getenv("API_KEY_ONESENDER"),
	}
	result, err := utils.CallAPI(http.MethodPost, url, &jsonData, headers, nil)
	if err != nil {
		utils.PrintLog("error [services][logics][sender][CallAPI] ", err)
		logrus.Error("error [services][logics][sender][CallAPI] ", err)
	}
	defer result.Body.Close()
}
