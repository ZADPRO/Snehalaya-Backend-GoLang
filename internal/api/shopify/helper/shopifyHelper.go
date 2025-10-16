package shopifyHelper

import (
	"net/http"

	"github.com/gin-gonic/gin"

)

func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   message,
	})
}
