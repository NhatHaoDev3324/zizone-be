package router

import (
	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/internal/middleware"
	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/internal/modules/user"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	r.LoadHTMLGlob("templates/*")

	r.GET("/", RootHandler)

	api := r.Group("/api/v1")
	{
		user.RegisterRoutes(api, db, rdb)
	}

	return r
}
