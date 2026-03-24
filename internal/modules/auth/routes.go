package auth

import (
	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/handler"
	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/repository"
	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.RouterGroup, db *gorm.DB, redis *redis.Client) {
	repo := repository.NewUserRepository(db, redis)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register-by-email", h.RegisterByEmail)
		userGroup.POST("/register-by-google", h.RegisterByGoogle)
		userGroup.POST("/login-by-email", h.LoginByEmail)
		userGroup.POST("/verify-otp", h.VerifyOTP)
		userGroup.GET("/", h.GetUsers)
		userGroup.GET("/:id", h.GetUserByID)
	}
}
