package logics

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user/constants"
	"user/domains/entities"
	"user/domains/models"
	"user/services/repositories"
	"user/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

type ITransactionService interface {
	// CreateTransaction(req models.Transaction) (*models.ResultToBuyer, error)
	CreateTransaction(req models.Transaction) (*models.ResultToBuyer, error)
	GetWebsocket(ctx *gin.Context)
	GetMelody(ctx *gin.Context)
	GetTransaction(req []entities.Transactions) (*[]models.Transaction, error)
	GetTransactionByID(req string) (*models.Transaction, error)
}

type TransactionService struct {
	transactionRepository repositories.ITransactionRepository
	useTransaction        bool
}

func InitTransactionService(transactionRepo repositories.ITransactionRepository) *TransactionService {
	if utils.IsNil(transactionRepo) {
		transactionRepo = repositories.InitTransactionRepository(nil, nil)
	}

	service := TransactionService{
		transactionRepository: transactionRepo,
		useTransaction:        false,
	}
	return &service
}

func (service *TransactionService) CreateTransaction(req models.Transaction) (*models.ResultToBuyer, error) {

	invoice_number := "INV"
	date := time.Now()
	timeFormat := date.Format("02-01-2006")
	rand_character := utils.RandomString(4)
	Invoice := fmt.Sprintf("%s-%s-%s", invoice_number, timeFormat, rand_character)
	rand.Seed(time.Now().UnixNano())
	rand_number := rand.Intn(99)

	TransactionToDB := entities.Transactions{}

	TransactionToDB.JenisLayanan = req.JenisLayanan
	TransactionToDB.Invoice_Number = Invoice
	TransactionToDB.Buyer_Sku_Code = req.Buyer_Sku_Code
	TransactionToDB.Nominal = req.Nominal
	TransactionToDB.Pembayaran = req.Pembayaran
	TransactionToDB.Nomor_Hp = req.Nomor_Hp
	TransactionToDB.Status = "Tunggu"
	TransactionToDB.Status_Pengisian = "Tunggu"
	TransactionToDB.Provider = req.Provider
	TransactionToDB.Id_Pelanggan = req.ID_Pelanggan
	TransactionToDB.Kode_Unik = rand_number
	TransactionToDB.Total = req.Nominal + rand_number

	// fetch user data
	transDataFirst, err := service.transactionRepository.CreateTransaction(&TransactionToDB)
	if err != nil {
		utils.PrintLog("error [services][logics][transaction][CreateTransaction] ", err)
		logrus.Error("error [services][logics][transaction][CreateTransaction] ", err)
		return nil, errors.New(constants.TransactionNotCreatedErr)
	}

	expiredTime := transDataFirst.CreatedAt.Add(10 * time.Minute)

	go service.CheckExpiredStatus(expiredTime, transDataFirst)

	ResultToBuyer := models.ResultToBuyer{}
	ResultToBuyer.Invoice_Number = transDataFirst.Invoice_Number
	ResultToBuyer.Unique_Code = transDataFirst.Kode_Unik
	ResultToBuyer.Total = transDataFirst.Total
	return &ResultToBuyer, nil
}

func (service *TransactionService) GetWebsocket(ctx *gin.Context) {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	var msg models.Message
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}
	// if err != nil {
	// 	utils.PrintLog("error [services][logics][transaction][GetWebsocket][Upgrade] ", err)
	// 	logrus.Error("error [services][logics][transaction][GetWebsocket][Upgrade] ", err)
	// }

	for {
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		fmt.Printf("Receive Message %s", msg)
		// conn.WriteMessage(t, msg)
		conn.SetWriteDeadline(time.Now().Add(60 * time.Minute))
		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello Client"))
		if err != nil {
			fmt.Printf("error sending message: %s\n", err.Error())
		}
		// if err != nil {
		// 	fmt.Println(err)
		// 	utils.PrintLog("error [services][logics][transaction][GetWebsocket][WriteMessage] ", err)
		// 	logrus.Error("error [services][logics][transaction][GetWebsocket][WriteMessage] ", err)
		// }
		conn.SetWriteDeadline(time.Time{})
	}
}

func (service *TransactionService) GetMelody(ctx *gin.Context) {
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	m.HandleRequest(ctx.Writer, ctx.Request)
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})
}

func (service *TransactionService) GetTransaction(req []entities.Transactions) (*[]models.Transaction, error) {

	var (
		transactions []models.Transaction
		transaction  models.Transaction
	)

	transactionData, err := service.transactionRepository.GetTransaction(req)
	if err != nil {
		logrus.Error("error [services][logics][transaction][GetTransaction] ", err)
		return nil, errors.New(constants.UserNotFoundErr)
	}

	for _, data := range transactionData {
		original := data.Nomor_Hp
		replacement := "XXX"
		startIndex := len(original) - 3
		endIndex := len(original)
		replaced := strings.Replace(original, original[startIndex:endIndex], replacement, 1)

		transaction.ID = int(data.ID)
		transaction.Status = data.Status
		transaction.JenisLayanan = data.JenisLayanan
		transaction.Nominal = data.Nominal

		transaction.Nomor_Hp = replaced
		transaction.Pembayaran = data.Pembayaran
		transaction.Provider = data.Provider
		transaction.CreatedAt = data.CreatedAt
		transaction.Buyer_Sku_Code = data.Buyer_Sku_Code
		transaction.Status_Pengisian = data.Status_Pengisian
		transaction.Kode_Unik = data.Kode_Unik
		transaction.Total = data.Total
		transaction.Serial_Number = data.Serial_Number
		transactions = append(transactions, transaction)
	}

	return &transactions, nil
}

func (service *TransactionService) GetTransactionByID(req string) (*models.Transaction, error) {
	var (
		transaction models.Transaction
	)
	ID, err := strconv.Atoi(req)
	if err != nil {
		logrus.Error("error [services][logics][transaction][strconv Atoi] ", err)
		return nil, errors.New(constants.UserNotFoundErr)
	}

	transactionData, err := service.transactionRepository.GetTransactionByID(ID)
	if err != nil {
		logrus.Error("error [services][logics][transaction][GetTransactionByID] ", err)
		return nil, errors.New(constants.UserNotFoundErr)
	}

	original := transactionData.Nomor_Hp
	replacement := "XXX"
	startIndex := len(original) - 3
	endIndex := len(original)
	replaced := strings.Replace(original, original[startIndex:endIndex], replacement, 1)

	transaction.ID = int(transactionData.ID)
	transaction.JenisLayanan = transactionData.JenisLayanan
	transaction.Provider = transactionData.Provider
	transaction.ID_Pelanggan = transactionData.Id_Pelanggan
	transaction.CreatedAt = transactionData.CreatedAt
	transaction.Kode_Unik = transactionData.Kode_Unik
	transaction.Nomor_Hp = replaced
	transaction.Pembayaran = transactionData.Pembayaran
	transaction.Status = transactionData.Status
	transaction.Nominal = transactionData.Nominal
	transaction.Status_Pengisian = transactionData.Status_Pengisian
	transaction.Total = transactionData.Total
	transaction.Serial_Number = transactionData.Serial_Number

	return &transaction, nil
}
