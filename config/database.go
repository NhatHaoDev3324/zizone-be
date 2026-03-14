package config

import (
	"fmt"
	"log"
	"os"

	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/internal/modules/user/model"

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
		log.Fatal("❌❌❌ Failed to connect database: ", err)
	}

	db.AutoMigrate(&model.User{})

	fmt.Println("✅✅✅ Connected to PostgreSQL successfully!")
	return db
}
