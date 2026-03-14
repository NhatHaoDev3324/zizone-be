package main

import (
	"fmt"
	"template/config"
	"template/internal/router"
)

func main() {
	config.LoadConfig()

	db := config.ConnectDB()
	rdb := config.ConnectRedis()

	r := router.NewRouter(db, rdb)

	fmt.Println("🚀🚀🚀 Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("❌❌❌ Server failed to start:", err)
	}
}
