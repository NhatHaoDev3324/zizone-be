package router

import (
	"template/internal/middleware"
	"template/internal/modules/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api/v1")
	{
		user.RegisterRoutes(api, db)
	}

	return r
}
