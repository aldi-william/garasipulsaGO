package logics

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"user/constants"
	"user/domains/entities"
	"user/domains/models"
	"user/services/repositories"
	"user/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ITransactionService interface {
	CreateTransaction(req models.Transaction) (*models.ResultToBuyer, error)
	CreateTransactionPLN(req models.TransactionPLN) (*models.ResultToBuyer, error)
	GetWebsocket(ctx *gin.Context)
	GetTransaction(req []entities.Transactions) (*[]models.Transaction, error)
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
	timeFormat := date.Format("02/01/2006")
	rand_character := utils.RandomString(4)
	Invoice := fmt.Sprintf("%s-%s-%s", invoice_number, timeFormat, rand_character)
	rand_number := rand.Intn(99)
	TransactionToDB := entities.Transactions{}

	TransactionToDB.JenisLayanan = req.JenisLayanan
	TransactionToDB.Invoice_Number = Invoice
	TransactionToDB.Buyer_Sku_Code = req.Buyer_Sku_Code
	TransactionToDB.Nominal = req.Nominal
	TransactionToDB.Pembayaran = req.Pembayaran
	TransactionToDB.Nomor_Hp = req.Nomor_Hp
	TransactionToDB.Status = "Tunggu"
	TransactionToDB.Provider = req.Provider
	TransactionToDB.Kode_Unik = rand_number
	TransactionToDB.Total = req.Nominal + rand_number

	// fetch user data
	transData, err := service.transactionRepository.CreateTransaction(&TransactionToDB)
	if err != nil {
		utils.PrintLog("error [services][logics][transaction][CreateTransaction] ", err)
		logrus.Error("error [services][logics][transaction][CreateTransaction] ", err)
		return nil, errors.New(constants.TransactionNotCreatedErr)
	}

	ResultToBuyer := models.ResultToBuyer{}

	ResultToBuyer.Invoice_Number = transData.Invoice_Number
	ResultToBuyer.Unique_Code = transData.Kode_Unik
	ResultToBuyer.Total = transData.Total

	return &ResultToBuyer, nil
}

func (service *TransactionService) CreateTransactionPLN(req models.TransactionPLN) (*models.ResultToBuyer, error) {

	invoice_number := "INV"
	date := time.Now()
	timeFormat := date.Format("02/01/2006")
	rand_character := utils.RandomString(4)
	Invoice := fmt.Sprintf("%s-%s-%s", invoice_number, timeFormat, rand_character)

	TransactionToDB := entities.TransactionsPLN{}

	TransactionToDB.Transactions.JenisLayanan = req.Transaction.JenisLayanan
	TransactionToDB.Transactions.Invoice_Number = Invoice
	TransactionToDB.Transactions.Buyer_Sku_Code = req.Transaction.Buyer_Sku_Code
	TransactionToDB.Transactions.Nominal = req.Transaction.Nominal
	TransactionToDB.Transactions.Pembayaran = req.Transaction.Pembayaran
	TransactionToDB.Transactions.Nomor_Hp = req.Transaction.Nomor_Hp
	TransactionToDB.Transactions.Status = "Tunggu"
	TransactionToDB.Transactions.Provider = req.Transaction.Provider
	TransactionToDB.Id_Pelanggan = req.ID_Pelanggan

	// fetch user data
	transData, err := service.transactionRepository.CreateTransactionPLN(&TransactionToDB)
	if err != nil {
		utils.PrintLog("error [services][logics][transaction][CreateTransaction] ", err)
		logrus.Error("error [services][logics][transaction][CreateTransaction] ", err)
		return nil, errors.New(constants.TransactionNotCreatedErr)
	}

	ResultToBuyer := models.ResultToBuyer{}

	ResultToBuyer.Invoice_Number = transData.Invoice_Number
	ResultToBuyer.Unique_Code = transData.Kode_Unik
	ResultToBuyer.Total = transData.Total

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
		transaction.Status = data.Status
		transaction.JenisLayanan = data.JenisLayanan
		transaction.Nominal = data.Nominal
		transaction.Nomor_Hp = data.Nomor_Hp
		transaction.Pembayaran = data.Pembayaran
		transaction.Provider = data.Provider
		transaction.CreatedAt = data.CreatedAt
		transactions = append(transactions, transaction)
	}

	return &transactions, nil
}
