package middleware

import (
	"net/http"
	"strings"

	"github.com/NhatHaoDev3324/zizone-be/pkg/response"
	"github.com/NhatHaoDev3324/zizone-be/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(ctx, http.StatusUnauthorized, "Token not found")
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Fail(ctx, http.StatusUnauthorized, "Token format is invalid")
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseAccessToken(tokenString)
		if err != nil || claims == nil {
			response.Fail(ctx, http.StatusUnauthorized, "Token is expired or invalid")
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.ID)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
