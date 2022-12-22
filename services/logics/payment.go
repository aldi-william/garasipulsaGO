package logics

import (
	"user/domains/models"
	"user/services/repositories"
	"user/utils"

	"github.com/sirupsen/logrus"
)

type IPaymentService interface {
	CallBackFromMoota(models.MootaCallback)
}

type PaymentService struct {
	paymentRepository     repositories.IPaymentRepository
	transactionRepository repositories.ITransactionRepository
	useTransaction        bool
}

func InitPaymentService(paymentRepo repositories.IPaymentRepository, transactionRepo repositories.ITransactionRepository) *PaymentService {
	if utils.IsNil(paymentRepo) {
		paymentRepo = repositories.InitPaymentRepository(nil, nil)
	}

	service := PaymentService{
		paymentRepository:     paymentRepo,
		transactionRepository: transactionRepo,
		useTransaction:        false,
	}
	return &service
}

func (service *PaymentService) CallBackFromMoota(req models.MootaCallback) {

	getData, err := service.transactionRepository.GetTransactionByTotal(req.Amount)
	if err != nil {
		utils.PrintLog("error [services][logics][payment][gorm get transactional by total] ", err)
		logrus.Error("error [services][logics][payment][gorm get transactional by total] ", err)
	}

	if getData.Total == req.Amount {
		utils.PrintLog("success [services][logics][payment][Total dan Amount Sama]", "SUKSES")
	}

}
