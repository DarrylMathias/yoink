package database

import (
	"fmt"
	"yoink/models"
	"yoink/utils/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewDatabase(env *env.Env){
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DBHost, env.DBPort, env.DBUser, env.DBPassword, env.DBName)
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		panic(fmt.Errorf("fatal error database connection: %w", err))
	}
	fmt.Println("Database connection success")
	database.AutoMigrate(&models.Page{})
	DB = database
}