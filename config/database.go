package config

import (
	"fmt"
	"os"

	"github.com/NhatHaoDev3324/goAuth/factory"
	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		factory.LogError("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(&model.User{})

	factory.LogSuccess("Connected to PostgreSQL successfully!")
	return db
}
