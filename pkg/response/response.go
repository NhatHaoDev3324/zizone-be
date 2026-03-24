package response

import (
	"net/http"

	"github.com/NhatHaoDev3324/goAuth/constant"
	"github.com/gin-gonic/gin"
)

func SuccessWithToken(ctx *gin.Context, message string, token string) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"token":   token,
	})
}

func SuccessWithData(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SuccessNoData(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
}

func Fail(ctx *gin.Context, status constant.FailStatus, message string) {
	ctx.JSON(int(status), gin.H{
		"success": false,
		"message": message,
	})
}
