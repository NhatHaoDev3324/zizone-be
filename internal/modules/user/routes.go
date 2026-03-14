package user

import (
	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/internal/modules/user/handler"
	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/internal/modules/user/repository"
	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/internal/modules/user/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	repo := repository.NewUserRepository(db, rdb)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", h.Register)
		userGroup.GET("/", h.GetUsers)
		userGroup.GET("/:id", h.GetUserByID)
	}
}
