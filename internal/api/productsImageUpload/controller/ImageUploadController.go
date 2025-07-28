package imageUploadController

import (
	"net/http"
	"strconv"

	imageUploadService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/productsImageUpload/service"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func CreateUploadURLHandler(c *gin.Context) {
	fileName := c.Param("filename")
	expireStr := c.Param("expireMinsDuration")

	expireMins, err := strconv.Atoi(expireStr)

	log.Info("Create Upload URL Handler - Begins ===> \n\n")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry"})
		log.Error("Invalid or Expire Token")
		return
	}

	uploadURL, fileURL, err := imageUploadService.CreateUploadURL(fileName, expireMins)
	log.Info("\n\nUpload URL -->", uploadURL)
	if err != nil {
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

	expireMins, err := strconv.Atoi(expireStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry"})
		return
	}

	fileURL, err := imageUploadService.GetFileURL(fileName, expireMins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate file URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fileUrl": fileURL})
}

func GetEnvVariables(c *gin.Context) {
	envVars := imageUploadService.FetchAllEnvVariables()
	c.JSON(http.StatusOK, gin.H{
		"env": envVars,
	})
}
