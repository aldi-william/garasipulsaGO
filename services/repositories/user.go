package repositories

import (
	"database/sql"
	"user/db"
	"user/domains/entities"
	"user/domains/models"
	"user/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	connORM *gorm.DB
	connDB  *sql.DB
}

type IUserRepository interface {
	SetConnection(connORM *gorm.DB)
	GetUserByID(id uint) (*entities.User, error)
	CreateUser(req models.Users) (*entities.Users, error)
	FindUserByEmail(email string) (*entities.Users, error)
}

func InitUserRepository(connORM *gorm.DB, connDB *sql.DB) *UserRepository {
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

// set connection
func (repo *UserRepository) SetConnection(connORM *gorm.DB) {
	repo.connORM = connORM
}

func (repo *UserRepository) GetUserByID(id uint) (*entities.User, error) {
	var result *entities.User
	err := repo.connORM.Table("USER").First(&result, id).Error
	if err != nil {
		logrus.Error("error [services][repositories][user][gorm read] ", err)
		return nil, err
	}
	return result, nil
}

func (repo *UserRepository) CreateUser(req models.Users) (*entities.Users, error) {
	user := entities.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
	err := repo.connORM.Create(&user).Error
	if err != nil {
		logrus.Error("error [services][repositories][user][gorm create] ", err)
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindUserByEmail(email string) (*entities.Users, error) {
	var user entities.Users
	err := repo.connORM.Where("email = ?", email).Find(&user).Error
	if err != nil {
		logrus.Error("error [services][repositories][user][gorm FindUserByEmail] ", err)
		return nil, err
	}
	return &user, nil
}
