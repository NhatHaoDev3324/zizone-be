package auth

import (
	"github.com/NhatHaoDev3324/zizone-be/internal/middleware"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/handler"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/repository"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.RouterGroup, db *gorm.DB, redis *redis.Client) {
	repo := repository.NewUserRepository(db, redis)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	authGroup := r.Group("/auth")
	{
		authGroup.GET("/profile", middleware.AuthMiddleware(), h.GetProfile)
		authGroup.POST("/register-by-email", h.RegisterByEmail)
		authGroup.POST("/register-by-google", h.RegisterByGoogle)
		authGroup.POST("/login-by-email", h.LoginByEmail)
		authGroup.POST("/verify-otp", h.VerifyOTP)
		authGroup.POST("/forgot-password", h.ForgotPassword)
		authGroup.POST("/verify-otp-forgot-password", h.VerifyOTPForgotPassword)
		authGroup.POST("/reset-password", middleware.AuthMiddleware(), h.ResetPassword)
	}
}
