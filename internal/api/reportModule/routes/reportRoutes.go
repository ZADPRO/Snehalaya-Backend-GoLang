package reportRoutes

import (
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

func ReportRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/reports")

	route.POST("/productReports", accesstoken.JWTMiddleware())
}
