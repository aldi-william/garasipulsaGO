package repositories

import (
	"database/sql"
	"user/db"
	"user/utils"

	"gorm.io/gorm"
)

type SenderRepository struct {
	connORM *gorm.DB
	connDB  *sql.DB
}

type ISenderRepository interface {
}

func InitSenderRepository(connORM *gorm.DB, connDB *sql.DB) *UserRepository {
	if utils.IsNil(connORM) {
		connORM = db.DBORM
	}

	if utils.IsNil(connDB) {
		connDB = db.DB
	}

	return &UserRepository{
		connORM: connORM,
		connDB:  connDB,
	}
}
