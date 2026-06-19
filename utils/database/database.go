package database

import (
	"fmt"
	"log"
	"os"
	"yoink/models"
	"yoink/utils/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func NewDatabase(env *env.Env){
	var dsn string
	if env.Application == "dev"{
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DBHost, env.DBPort, env.DBUser, env.DBPassword, env.DBName)
	}else{
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full sslrootcert=%s", env.DBHost, env.DBPort, env.DBUser, env.DBPassword, env.DBName, env.DBSSLRootCert)
	}
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
			LogLevel: logger.Error,
		},
	),
	})
	if err != nil{
		panic(fmt.Errorf("fatal error database connection: %w", err))
	}
	if env.Application == "dev"{
		fmt.Println("Local database connection success")
	}else{
		fmt.Println("RDS connection success")
	}
	database.AutoMigrate(&models.Page{})
	DB = database
}