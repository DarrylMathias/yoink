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

func NewDatabase(){
	var dsn string
	if env.ConfigValue.Application == "dev"{
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.EnvValue.DBHost, env.EnvValue.DBPort, env.EnvValue.DBUser, env.EnvValue.DBPassword, env.EnvValue.DBName)
	}else{
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full sslrootcert=%s", env.EnvValue.DBHost, env.EnvValue.DBPort, env.EnvValue.DBUser, env.EnvValue.DBPassword, env.EnvValue.DBName, env.EnvValue.DBSSLRootCert)
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
	if env.ConfigValue.Application == "dev"{
		fmt.Println("Local database connection success")
	}else{
		fmt.Println("RDS connection success")
	}
	database.AutoMigrate(&models.Page{})
	database.AutoMigrate(&models.Posting{})
	database.AutoMigrate(&models.Term{})
	DB = database
}