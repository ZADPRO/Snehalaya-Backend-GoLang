package imageUploadController

import (
	"net/http"
	"strconv"

	imageUploadService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/productsImageUpload/service"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

type UploadRequest struct {
	FileName   string `json:"fileName" binding:"required"`
	ExpireMins int    `json:"expireMins"`
}

func CreateUploadURLHandler(c *gin.Context) {
	var req UploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Invalid request body | Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set default expiry if not provided or invalid
	expireMins := req.ExpireMins
	if expireMins <= 0 {
		expireMins = 5
	}

	log.Infof("CreateUploadURLHandler called | fileName: %s | expireMins: %d", req.FileName, expireMins)

	uploadURL, fileURL, err := imageUploadService.CreateUploadURL(req.FileName, expireMins)
	if err != nil {
		log.Errorf("Failed to create presigned PUT URL | fileName: %s | expireMins: %d | Error: %+v", req.FileName, expireMins, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate upload URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uploadUrl": uploadURL,
		"fileUrl":   fileURL,
	})
}

func GetFileURLHandler(c *gin.Context) {
	fileName := c.Param("filename")
	expireStr := c.Param("expireMins")

	log.Infof("GetFileURLHandler called | fileName: %s | expireMins: %s", fileName, expireStr)

	expireMins, err := strconv.Atoi(expireStr)
	if err != nil {
		log.Errorf("Invalid expiry string: %s | Error: %v", expireStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry"})
		return
	}

	fileURL, err := imageUploadService.GetFileURL(fileName, expireMins)
	if err != nil {
		log.Errorf("Failed to generate presigned GET URL for file %s | Error: %v", fileName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate file URL"})
		return
	}

	log.Infof("File URL generated successfully | fileURL: %s", fileURL)
	c.JSON(http.StatusOK, gin.H{"fileUrl": fileURL})
}

func GetEnvVariables(c *gin.Context) {
	log.Info("GetEnvVariables called")
	envVars := imageUploadService.FetchAllEnvVariables()
	c.JSON(http.StatusOK, gin.H{
		"env": envVars,
	})
}
