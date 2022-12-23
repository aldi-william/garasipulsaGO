package logics

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"user/domains/models"
	"user/services/repositories"
	"user/utils"

	"github.com/sirupsen/logrus"
)

type IPaymentService interface {
	CallBackFromMoota(req []models.MootaCallback) (*models.ResultDigiFlazzData, error)
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

	if utils.IsNil(transactionRepo) {
		transactionRepo = repositories.InitTransactionRepository(nil, nil)
	}

	service := PaymentService{
		paymentRepository:     paymentRepo,
		transactionRepository: transactionRepo,
		useTransaction:        false,
	}
	return &service
}

func (service *PaymentService) CallBackFromMoota(req []models.MootaCallback) (*models.ResultDigiFlazzData, error) {

	var (
		amount           = []int{}
		res              = models.ResultDigiFlazzData{}
		MODE_DEVELOPMENT = os.Getenv("MODE_DEVELOPMENT")
		Testing          bool
	)
	// menentukan mode development atau production
	if MODE_DEVELOPMENT == "PRODUCTION" {
		Testing = false
	} else {
		Testing = true
	}

	for _, data := range req {
		amount = append(amount, data.Amount)
	}

	date := time.Now()
	format_date := date.Format("2006-01-02")

	getDataAllStatusTungguAndToday, err := service.transactionRepository.GetTransactionByStatusAndToday("Tunggu", format_date)
	if err != nil {
		utils.PrintLog("error [services][logics][payment][gorm get transactional by Status] ", err)
		logrus.Error("error [services][logics][payment][gorm get transactional by Status] ", err)
	}

	username := os.Getenv("USERNAME_BELI_PULSA")
	apikey := os.Getenv("API_KEY_BELI_PULSA")

	jsonData := models.TransactionDIGIFLAZZ{}
	jsonData.Command = os.Getenv("CMD_BELI_PULSA")
	jsonData.Username = os.Getenv("USERNAME_BELI_PULSA")

	for _, data := range amount {
		if getDataAllStatusTungguAndToday.Total == data {
			utils.PrintLog("success [services][logics][payment][getDataAllStatusTungguAndToday == Amount of all Mutasi] ", data)
			getData, err := service.transactionRepository.GetTransactionByTotal(data)
			if err != nil {
				utils.PrintLog("error [services][logics][payment][gorm get transactional by total] ", err)
				logrus.Error("error [services][logics][payment][gorm get transactional by total] ", err)
			}
			jsonData.Customer_NO = strconv.Itoa(getData.Nomor_Hp)
			jsonData.Buyer_SKU_Code = getData.Buyer_Sku_Code
			jsonData.Ref_ID = getData.Invoice_Number
			ref_id := getData.Invoice_Number
			sign := md5.Sum([]byte(username + apikey + ref_id))
			pass := fmt.Sprintf("%x", sign)
			jsonData.Sign = pass
			jsonData.Testing = Testing
			result, err := utils.CallAPI(http.MethodPost, os.Getenv("URL_BELI_PULSA"), &jsonData, nil, nil)
			if err != nil {
				utils.PrintLog("error [services][logics][payment][CallAPI] ", err)
				logrus.Error("error [services][logics][payment][CallAPI] ", err)
			}
			defer result.Body.Close()
			bytes, err := io.ReadAll(result.Body)
			if err != nil {
				utils.PrintLog("error [services][logics][payment][ReadAll Looping CallAPI] ", err)
				logrus.Error("error [services][logics][payment][ReadAll Looping CallAPI] ", err)
			}
			err = json.Unmarshal(bytes, &res)
			if err != nil {
				utils.PrintLog("error [services][logics][payment][Unmarshal Looping CallAPI] ", err)
				logrus.Error("error [services][logics][payment][Unmarshal Looping CallAPI] ", err)
			}

		} else {
			fmt.Println("Tidak Ada Transaksi")
		}
	}

	return &res, nil
}
