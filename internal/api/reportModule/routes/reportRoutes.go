package reportRoutes

import (
	reportController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

func ReportRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/reports")

	route.POST("/productReports", accesstoken.JWTMiddleware(), reportController.GetAllProductsReportController())
	route.POST("/productReportsDownload", accesstoken.JWTMiddleware(), reportController.GetAllProductsReportController())

}
