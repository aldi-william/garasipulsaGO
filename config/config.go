package config

import (
	"os"
	"user/connection"
	"user/db"
	"user/server"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// StartApp is used to Start the Application
func StartApp() {
	setupConfig()
	server.RegisterRoutes()
	server.ListenAndServe()
}

// setupConfig is used to load environment variable
func setupConfig() {
	envVariable := godotenv.Load(".env")
	if envVariable != nil {
		logrus.Fatalln("Error loading .env file")
	}

	initConfig()
}

func initConfig() {
	var (
		environment string = os.Getenv("ENVIRONMENT")
	)

	// setup connection DB
	err := db.InitDB()
	if err != nil {
		logrus.Fatal(err)
	}

	// setup connection http
	connection.InitClient()

	server.AllowCors(environment)
}

//endregion functions
