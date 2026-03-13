package user

import (
	"template/internal/modules/user/handler"
	"template/internal/modules/user/repository"
	"template/internal/modules/user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", h.Register)
		userGroup.GET("/", h.GetUsers)
	}
}
