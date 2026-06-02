package main

import (
	"github.com/NhatHaoDev3324/zizone-be/config"
	"github.com/NhatHaoDev3324/zizone-be/internal/router"
	"github.com/NhatHaoDev3324/zizone-be/pkg/log"
	"github.com/NhatHaoDev3324/zizone-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.LogInfo(".env not found, fallback to .env.local")
		godotenv.Load(".env.local")
	} else {
		log.LogSuccess("Environment loaded from .env")
	}

	db := config.ConnectDB()
	redis := config.ConnectRedis()
	config.InitCloudinary()

	utils.NewMailService(5)

	gin.SetMode(gin.ReleaseMode)

	r := router.NewRouter(db, redis)
	r.SetTrustedProxies([]string{"nil"})

	log.LogSuccess("Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.LogError("Server failed to start: " + err.Error())
	}
}
