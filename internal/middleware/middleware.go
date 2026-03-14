package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/NhatHaoDev3324/go-gin-gorm-postgres-template/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Token không tồn tại"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Sai định dạng token"})
			ctx.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := utils.ParseAccessToken(tokenString)
		if err != nil || claims == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Token hết hạn hoặc không hợp lệ"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", fmt.Sprintf("%v", claims.ID))
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
