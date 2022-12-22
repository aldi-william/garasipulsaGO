package repositories

import (
	"database/sql"
	"user/db"
	"user/utils"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	connORM *gorm.DB
	connDB  *sql.DB
}

type IPaymentRepository interface {
}

func InitPaymentRepository(connORM *gorm.DB, connDB *sql.DB) *TransactionRepository {
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
