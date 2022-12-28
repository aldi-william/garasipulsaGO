package db

//region imports
import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//endregion imports

const (
	driverName = "mysql"
)

// region variables
var (
	DBORM *gorm.DB
	DB    *sql.DB
)

//endregion variables

func InitDB() error {
	var (
		dbHost     string = os.Getenv("DB_HOST")
		dbUser     string = os.Getenv("DB_USER")
		dbPassword string = os.Getenv("DB_PASSWORD")
		dbPort     string = os.Getenv("DB_PORT")
		dbName     string = os.Getenv("DB_NAME")
		err        error
	)

	// setup connection DB postgres  dbUser:dbPassword@tcp(dbHost:dbPort)/dbName
	// db, err := sql.Open("mysql", ":Mysql123!@tcp(localhost:3306)/user")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	DB, err = sql.Open(driverName, dsn)
	if err != nil {
		return errors.Wrap(err, "unable to open db connection")
	}

	err = DB.Ping()
	if err != nil {
		return errors.Wrap(err, "failed ping db")
	}

	DBORM, err = gorm.Open(
		mysql.New(mysql.Config{
			Conn: DB,
		}), &gorm.Config{
			SkipDefaultTransaction: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		return errors.Wrap(err, "unable to init gorm")
	}

	err = autoMigrate(DB)
	if err != nil {
		return errors.Wrap(err, "failed auto migration")
	}
	return nil
}
