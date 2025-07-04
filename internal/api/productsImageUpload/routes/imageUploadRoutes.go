package imageUploadRoutes

import (
	imageUploadController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/productsImageUpload/controller"
	"github.com/gin-gonic/gin"
)

func ImageUploadRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/imageUpload")
	route.GET("/productImages/:filename/:expireMinsDuration", imageUploadController.CreateUploadURLHandler)
	route.GET("/getProductImage/:filename/:expireMins", imageUploadController.GetFileURLHandler)
}
