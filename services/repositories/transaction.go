package repositories

import (
	"database/sql"
	"user/db"
	"user/domains/entities"
	"user/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	connORM *gorm.DB
	connDB  *sql.DB
}

type ITransactionRepository interface {
	PutTransaction(req *entities.Transactions) (string, error)
	CreateTransaction(req *entities.Transactions) (*entities.Transactions, error)
	GetTransaction(req []entities.Transactions) ([]entities.Transactions, error)
	GetTransactionByGagal(req []entities.Transactions) ([]entities.Transactions, error)
	GetTransactionByInvoice(invoice string) (*entities.Transactions, error)
	GetTransactionByTotal(total int) (*entities.Transactions, error)
	GetTransactionByStatusAndToday(status string, date string) (*[]entities.Transactions, error)
	GetTransactionByID(id int) (*entities.Transactions, error)
	UpdateTransactionByInvoiceNumber(trans *entities.Transactions) error
}

func InitTransactionRepository(connORM *gorm.DB, connDB *sql.DB) *TransactionRepository {
	if utils.IsNil(connORM) {
		connORM = db.DBORM
	}

	if utils.IsNil(connDB) {
		connDB = db.DB
	}

	return &TransactionRepository{
		connORM: connORM,
		connDB:  connDB,
	}
}

func (transactionRepo *TransactionRepository) CreateTransaction(req *entities.Transactions) (*entities.Transactions, error) {
	err := transactionRepo.connORM.Create(&req).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm create Transaction PLN] ", err)
		logrus.Error("error [services][repositories][transaction][gorm create Transaction PLN] ", err)
		return nil, err
	}
	return req, nil
}

func (transactionRepo *TransactionRepository) GetTransaction(req []entities.Transactions) ([]entities.Transactions, error) {
	err := transactionRepo.connORM.Limit(20).Order("created_at desc").Find(&req).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm get] ", err)
		logrus.Error("error [services][repositories][transaction][gorm get] ", err)
		return nil, err
	}
	return req, nil
}

func (transactionRepo *TransactionRepository) GetTransactionByGagal(req []entities.Transactions) ([]entities.Transactions, error) {
	err := transactionRepo.connORM.Order("created_at desc").Where("status_pengisian = ?", "Gagal").Find(&req).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm get] ", err)
		logrus.Error("error [services][repositories][transaction][gorm get] ", err)
		return nil, err
	}
	return req, nil
}

func (transactionRepo *TransactionRepository) GetTransactionByID(id int) (*entities.Transactions, error) {
	var (
		trans entities.Transactions
	)
	err := transactionRepo.connORM.Where("id = ?", id).First(&trans).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm get] ", err)
		logrus.Error("error [services][repositories][transaction][gorm get] ", err)
		return nil, err
	}
	return &trans, nil
}

func (transactionRepo *TransactionRepository) GetTransactionByInvoice(invoice string) (*entities.Transactions, error) {
	var (
		trans entities.Transactions
	)
	err := transactionRepo.connORM.Where("invoice_number = ?", invoice).First(&trans).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm get] ", err)
		logrus.Error("error [services][repositories][transaction][gorm get] ", err)
		return nil, err
	}

	return &trans, nil
}

func (transactionRepo *TransactionRepository) GetTransactionByTotal(total int) (*entities.Transactions, error) {
	var (
		trans entities.Transactions
	)
	err := transactionRepo.connORM.Where("total = ?", total).First(&trans).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm get] ", err)
		logrus.Error("error [services][repositories][transaction][gorm get] ", err)
		return nil, err
	}

	return &trans, nil
}

func (transactionRepo *TransactionRepository) GetTransactionByStatusAndToday(status string, date string) (*[]entities.Transactions, error) {
	var (
		trans *[]entities.Transactions
	)
	err := transactionRepo.connORM.Where("status = ? AND created_at > ?", status, date).Find(&trans).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm get GET TransactionByStatusAndToday] ", err)
		logrus.Error("error [services][repositories][transaction][gorm get GET TransactionByStatusAndToday] ", err)
		return nil, err
	}

	return trans, nil
}

func (transactionRepo *TransactionRepository) UpdateTransactionByInvoiceNumber(trans *entities.Transactions) error {
	var (
		transaction entities.Transactions
	)

	err := transactionRepo.connORM.Model(&transaction).Where("invoice_number = ?", trans.Invoice_Number).Updates(entities.Transactions{
		Status:           trans.Status,
		Serial_Number:    trans.Serial_Number,
		Status_Pengisian: trans.Status_Pengisian,
	}).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm Update UpdateTransactionByInvoiceNumber] ", err)
		logrus.Error("error [services][repositories][transaction][gorm Update UpdateTransactionByInvoiceNumber] ", err)
		return err
	}
	return nil
}

func (transactionRepo *TransactionRepository) PutTransaction(req *entities.Transactions) (string, error) {

	err := transactionRepo.connORM.Model(&entities.Transactions{}).Where("invoice_number = ?", req.Invoice_Number).Update("Status", req.Status).Error
	if err != nil {
		utils.PrintLog("error [services][repositories][transaction][gorm Update UpdateTransactionByInvoiceNumber] ", err)
		logrus.Error("error [services][repositories][transaction][gorm Update UpdateTransactionByInvoiceNumber] ", err)
		return "gagal update transaction", err
	}

	return "sukses update transaction", nil
}
