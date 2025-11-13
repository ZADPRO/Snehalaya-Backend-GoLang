package bulkImageUploadRoutes

import (
	bulkImageUploadController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/bulkImageHandling/controller"
	"github.com/gin-gonic/gin"

)

func BulkImageUploadRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/bulkImageUpload")

	// Generate presigned URLs for bulk image upload
	route.POST("/generateUploadURL", bulkImageUploadController.GenerateBulkUploadURLHandler)

	// Optional: Generate presigned GET URL for temporary preview
	route.GET("/getImageViewURL/:filename/:expireMins", bulkImageUploadController.GetImageViewURLHandler)
}
