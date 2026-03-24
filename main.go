package main

import (
	"github.com/NhatHaoDev3324/goAuth/config"
	"github.com/NhatHaoDev3324/goAuth/factory"
	"github.com/NhatHaoDev3324/goAuth/internal/router"
	"github.com/NhatHaoDev3324/goAuth/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		factory.LogInfo(".env not found, fallback to .env.local")
		godotenv.Load(".env.local")
	} else {
		factory.LogSuccess("Environment loaded from .env")
	}

	db := config.ConnectDB()
	redis := config.ConnectRedis()

	utils.NewMailService(5)

	gin.SetMode(gin.ReleaseMode)

	r := router.NewRouter(db, redis)
	r.SetTrustedProxies([]string{"nil"})

	factory.LogSuccess("Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		factory.LogError("Server failed to start: " + err.Error())
	}
}
