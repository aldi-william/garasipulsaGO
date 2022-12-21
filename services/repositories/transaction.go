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
	CreateTransaction(req *entities.Transactions) (*entities.Transactions, error)
	CreateTransactionPLN(req *entities.TransactionsPLN) (*entities.TransactionsPLN, error)
	GetTransaction(req []entities.Transactions) ([]entities.Transactions, error)
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
		logrus.Error("error [services][repositories][transaction][gorm create] ", err)
		return nil, err
	}
	return req, nil
}

func (transactionRepo *TransactionRepository) CreateTransactionPLN(req *entities.TransactionsPLN) (*entities.TransactionsPLN, error) {
	err := transactionRepo.connORM.Create(&req).Error
	if err != nil {
		logrus.Error("error [services][repositories][transaction][gorm create] ", err)
		return nil, err
	}
	return req, nil
}

func (transactionRepo *TransactionRepository) GetTransaction(req []entities.Transactions) ([]entities.Transactions, error) {
	err := transactionRepo.connORM.Limit(10).Order("created_at desc").Find(&req).Error
	if err != nil {
		logrus.Error("error [services][repositories][transaction][gorm get] ", err)
		return nil, err
	}
	return req, nil
}
