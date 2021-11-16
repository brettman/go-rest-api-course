package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"os"
)

// NewDatabase - returns pointer to a new db object
func NewDatabase() (*gorm.DB, error){
	log.Info("Setting up new database connection")

	dbUserName := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable:= os.Getenv("DB_TABLE")
	dbPort:= os.Getenv("DB_PORT")
	sslMode:= os.Getenv("SSL_MODE")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUserName, dbTable, dbPassword, sslMode)
	db, err := gorm.Open("postgres", connectionString)
	if err !=nil{
		fmt.Printf("Error opening new database connection: %s\n", connectionString)
		return db, err
	}
	if err:= db.DB().Ping(); err !=nil{
		fmt.Println("Error pinging new database connection")
		return db, err
	}

	return db, nil
}