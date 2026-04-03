package router

import (
	"net/http"

	"github.com/NhatHaoDev3324/zizone-be/internal/middleware"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, redis *redis.Client) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ParseJWT())

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	api := r.Group("/api/v1")
	{
		auth.AuthRoutes(api, db, redis)
	}

	return r
}
