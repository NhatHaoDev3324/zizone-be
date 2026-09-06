package response

import (
	"net/http"

	"github.com/NhatHaoDev3324/zizone-be/internal/tdo"
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

func SuccessDataInfo(ctx *gin.Context, message string, data tdo.Profile) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SuccessWithMetaAndData(ctx *gin.Context, message string, meta tdo.Meta, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"meta":    meta,
		"data":    data,
	})
}

func Fail(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}
