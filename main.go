package main

import (
	"fmt"

	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/config"
	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Println("⚠️⚠️⚠️  .env.local not found, fallback to .env")
		godotenv.Load(".env")
	} else {
		fmt.Println("✅✅✅ Environment loaded from .env.local")
	}

	db := config.ConnectDB()
	rdb := config.ConnectRedis()

	gin.SetMode(gin.ReleaseMode)
	r := router.NewRouter(db, rdb)
	r.SetTrustedProxies([]string{"nil"})

	fmt.Println("🚀🚀🚀 Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("❌❌❌ Server failed to start:", err)
	}
}
