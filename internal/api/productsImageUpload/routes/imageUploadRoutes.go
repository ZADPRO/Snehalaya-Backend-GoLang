package imageUploadRoutes

import (
	imageUploadController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/productsImageUpload/controller"
	"github.com/gin-gonic/gin"
)

func ImageUploadRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/imageUpload")
	route.POST("/productImages", imageUploadController.CreateUploadURLHandler)
	route.GET("/getProductImage/:filename/:expireMins", imageUploadController.GetFileURLHandler)

	route.GET("/env", imageUploadController.GetEnvVariables)

	route.POST("/generateURL", imageUploadController.GetPresignedURL)

}
