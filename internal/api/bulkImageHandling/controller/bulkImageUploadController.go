package bulkImageUploadController

import (
	"net/http"
	"strconv"
	"strings"

	bulkImageUploadService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/bulkImageHandling/service"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

type BulkUploadRequest struct {
	FileNames  []string `json:"fileNames" binding:"required"` // multiple image names
	ExpireMins int      `json:"expireMins"`
}

// Generate presigned PUT + GET URLs for multiple files
func GenerateBulkUploadURLHandler(c *gin.Context) {
	var req BulkUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.FileNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file names provided"})
		return
	}

	expireMins := req.ExpireMins
	if expireMins <= 0 {
		expireMins = 15
	}

	results := []map[string]string{}

	for _, name := range req.FileNames {
		fileName := strings.ToUpper(name)
		uploadURL, viewURL, err := bulkImageUploadService.CreatePresignedURLs(fileName, expireMins)
		if err != nil {
			log.Errorf("Failed to create presigned URLs for %s: %v", fileName, err)
			continue
		}

		results = append(results, map[string]string{
			"fileName":  fileName,
			"uploadUrl": uploadURL,
			"viewUrl":   viewURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"results": results,
	})
}

func GetImageViewURLHandler(c *gin.Context) {
	fileName := c.Param("filename")
	expireStr := c.Param("expireMins")

	expireMins, err := strconv.Atoi(expireStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry"})
		return
	}

	url, err := bulkImageUploadService.GetImageViewURL(fileName, expireMins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate view URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"viewUrl": url,
	})
}
