package main

import (
	"fmt"

	"github.com/NhatHaoDev3324/GoTemplate/config"
	"github.com/NhatHaoDev3324/GoTemplate/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("📢 .env not found, fallback to .env.local")
		godotenv.Load(".env.local")
	} else {
		fmt.Println("✅ Environment loaded from .env")
	}

	db := config.ConnectDB()
	redis := config.ConnectRedis()

	gin.SetMode(gin.ReleaseMode)

	r := router.NewRouter(db, redis)
	r.SetTrustedProxies([]string{"nil"})

	fmt.Println("🚀 Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("❌ Server failed to start:", err)
	}
}
