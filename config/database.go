package config

import (
	"fmt"
	"os"

	authModel "github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/model"
	wordModel "github.com/NhatHaoDev3324/zizone-be/internal/modules/word/model"
	"github.com/NhatHaoDev3324/zizone-be/pkg/log"

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
		log.LogError("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(&authModel.User{}, &wordModel.Word{})

	log.LogSuccess("Connected to PostgreSQL successfully!")
	return db
}
