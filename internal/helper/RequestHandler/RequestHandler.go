package helper

import (
	"fmt"
	"net/http"

	hashapi "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/HashAPI"
	model "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func RequestHandler[T any](c *gin.Context) (*T, bool) {
	// EXTRACK TOKEN FROM CONTEXT
	tokenVal, exists := c.Get("token")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Token not found in context.",
		})
		return nil, false

	}
	// BIND ENCRYPTED BODY
	var encryptedData model.ReqVal
	if err := c.BindJSON(&encryptedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body : " + err.Error(),
		})
		return nil, false
	}
	if len(encryptedData.EncryptedData) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid encrypted data format",
		})
		return nil, false
	}

	// DECRYPTED DATA
	decryptedInterface, err := hashapi.Decrypt(encryptedData.EncryptedData, tokenVal.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Decryption failed: " + err.Error(),
		})
		return nil, false
	}

	// VALIDATE DECRYPTED SSTRUCURE
	mapData, ok := decryptedInterface.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Invalid decrypt format",
		})
		return nil, false
	}

	// PRINT DATA
	fmt.Println("\n\nDecoded Struct Data : ", mapData)

	var data T
	if err := mapstructure.Decode(mapData, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to decode decrypted data : " + err.Error(),
		})
		return nil, false
	}

	return &data, true
}
