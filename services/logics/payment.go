package logics

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"user/domains/entities"
	"user/domains/models"
	"user/onesender"
	"user/services/repositories"
	"user/utils"

	"github.com/sirupsen/logrus"
)

type IPaymentService interface {
	CallBackFromMoota(req []models.MootaCallback) (*models.ResultDigiFlazzData, error)
	CallBackFromDigiflazz(req models.DigiFlazzData) (*models.ResultDigiFlazzData, error)
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
		transaction      = entities.Transactions{}
	)

	onesender.ApiUrl = os.Getenv("API_URL_ONESENDER")
	onesender.ApiKey = os.Getenv("API_KEY_ONESENDER")
	// menentukan mode development atau production
	if MODE_DEVELOPMENT == "PRODUCTION" {
		Testing = false
	} else {
		Testing = true
	}

	for _, data := range req {
		intamount, _ := strconv.Atoi(data.Amount)
		amount = append(amount, intamount)
	}

	utils.PrintLogSukses("Amount", amount)

	date := time.Now()
	format_date := date.Format("2006-01-02")

	getDataAllStatusTungguAndToday, err := service.transactionRepository.GetTransactionByStatusAndToday("Tunggu", format_date)
	if err != nil {
		utils.PrintLog("error [services][logics][payment][gorm get transactional by Status] ", err)
		logrus.Error("error [services][logics][payment][gorm get transactional by Status] ", err)
		return nil, errors.New("data tidak ditemukan")
	}

	username := os.Getenv("USERNAME_BELI_PULSA")
	apikey := os.Getenv("API_KEY_BELI_PULSA")

	jsonData := models.TransactionDIGIFLAZZ{}
	jsonData.Command = os.Getenv("CMD_BELI_PULSA")
	jsonData.Username = os.Getenv("USERNAME_BELI_PULSA")

	for _, getData := range *getDataAllStatusTungguAndToday {
		for _, getDataAmount := range amount {
			if getData.Total == getDataAmount {
				utils.PrintLogSukses("success [services][logics][payment][getDataAllStatusTungguAndToday == Amount of all Mutasi] ", getDataAmount)
				getData, err := service.transactionRepository.GetTransactionByTotal(getDataAmount)
				if err != nil {
					utils.PrintLog("error [services][logics][payment][gorm get transactional by total] ", err)
					logrus.Error("error [services][logics][payment][gorm get transactional by total] ", err)
					return nil, errors.New("data tidak ditemukan")
				}
				if getData.Id_Pelanggan != "" {
					no_hp := fmt.Sprintf("0%s", getData.Nomor_Hp)
					jsonData.Customer_NO = no_hp
				} else {
					jsonData.Customer_NO = getData.Id_Pelanggan
				}

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

				if res.Data.Response_Code == "00" {
					transaction.Status = "Sukses"
					transaction.Status_Pengisian = "Sukses"
					transaction.Serial_Number = res.Data.Serial_Number
					msg := fmt.Sprintf("serial number anda adalah %s", transaction.Serial_Number)
					onesender.SendTextMessage(res.Data.Customer_No, msg)
				} else if res.Data.Response_Code == "01" {
					transaction.Status = "Sukses"
					transaction.Status_Pengisian = "Gagal"
					transaction.Serial_Number = ""
				} else if res.Data.Response_Code == "02" {
					transaction.Status = "Sukses"
					transaction.Status_Pengisian = "Gagal"
					transaction.Serial_Number = ""
				} else if res.Data.Response_Code == "03" {
					transaction.Status = "Sukses"
					transaction.Status_Pengisian = "Proses"
					transaction.Serial_Number = ""
				} else if res.Data.Response_Code == "99" {
					transaction.Status = "Sukses"
					transaction.Status_Pengisian = "Proses"
					transaction.Serial_Number = ""
				} else {
					transaction.Status = "Sukses"
					transaction.Status_Pengisian = "Gagal"
					transaction.Serial_Number = ""
				}

				transaction.Invoice_Number = res.Data.Ref_ID
				err = service.transactionRepository.UpdateTransactionByInvoiceNumber(&transaction)
				if err != nil {
					utils.PrintLog("error [services][logics][payment][Unmarshal Looping CallAPI] ", err)
					logrus.Error("error [services][logics][payment][Unmarshal Looping CallAPI] ", err)
				}
			} else {
				fmt.Println("Tidak Ada Transaksi")
			}
		}

	}

	return nil, nil
}

func (service *PaymentService) CallBackFromDigiflazz(req models.DigiFlazzData) (*models.ResultDigiFlazzData, error) {
	var (
		transaction = entities.Transactions{}
	)

	if req.Data.Response_Code == "00" {
		transaction.Status = "Sukses"
		transaction.Status_Pengisian = "Sukses"
		transaction.Serial_Number = req.Data.Serial_Number
	} else if req.Data.Response_Code == "01" {
		transaction.Status = "Sukses"
		transaction.Status_Pengisian = "Gagal"
		transaction.Serial_Number = ""
	} else if req.Data.Response_Code == "02" {
		transaction.Status = "Sukses"
		transaction.Status_Pengisian = "Gagal"
		transaction.Serial_Number = ""
	} else if req.Data.Response_Code == "03" {
		transaction.Status = "Sukses"
		transaction.Status_Pengisian = "Proses"
		transaction.Serial_Number = ""
	} else if req.Data.Response_Code == "99" {
		transaction.Status = "Sukses"
		transaction.Status_Pengisian = "Proses"
		transaction.Serial_Number = ""
	} else {
		transaction.Status = "Sukses"
		transaction.Status_Pengisian = "Gagal"
		transaction.Serial_Number = ""
	}

	transaction.Invoice_Number = req.Data.Ref_ID

	err := service.transactionRepository.UpdateTransactionByInvoiceNumber(&transaction)
	if err != nil {
		utils.PrintLog("error [services][logics][payment][Unmarshal Looping CallAPI] ", err)
		logrus.Error("error [services][logics][payment][Unmarshal Looping CallAPI] ", err)
		return nil, errors.New("gagal update transaksi")
	}
	utils.PrintLogSukses("SUCCESS [services][logics][payment][CALLBACK FROM Digiflazz] ", transaction)
	return nil, nil
}
