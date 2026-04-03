package middleware

import (
	"net/http"
	"strings"

	"github.com/NhatHaoDev3324/zizone-be/pkg/response"
	"github.com/NhatHaoDev3324/zizone-be/utils"

	"github.com/gin-gonic/gin"
)

func ParseJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")

			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString := parts[1]

				claims, err := utils.ParseAccessToken(tokenString)
				if err == nil && claims != nil {
					ctx.Set("userID", claims.ID)
					ctx.Set("role", claims.Role)
				}

			}
		}
		ctx.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userID, exists := ctx.Get("userID")
		if !exists || userID == nil {
			response.Fail(ctx, http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		roleRaw, exists := ctx.Get("role")
		if !exists {
			response.Fail(ctx, http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		role := roleRaw.(string)

		for _, r := range roles {
			if role == r {
				ctx.Next()
				return
			}
		}

		response.Fail(ctx, http.StatusForbidden, "Forbidden")
		ctx.Abort()
	}
}
