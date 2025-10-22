package shopifyController

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	shopifyHelper "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/shopify/helper"
	shopifyService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/shopify/service"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/gin-gonic/gin"

)

func GetShopifyProducts(ctx *gin.Context) {
	products, err := shopifyService.GetAllProducts()
	if err != nil {
		shopifyHelper.ErrorResponse(ctx, err.Error())
		return
	}
	shopifyHelper.SuccessResponse(ctx, products)
}

func CreateShopifyProduct(ctx *gin.Context) {
	var product goshopify.Product

	// Bind JSON payload from request
	if err := ctx.ShouldBindJSON(&product); err != nil {
		shopifyHelper.ErrorResponse(ctx, "Invalid payload: "+err.Error())
		return
	}

	createdProduct, err := shopifyService.CreateProduct(product)
	if err != nil {
		shopifyHelper.ErrorResponse(ctx, err.Error())
		return
	}

	shopifyHelper.SuccessResponse(ctx, createdProduct)
}

func OrderCreationWebhook() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("\n\n\n\nOrder creation web hook called ->>>> \n")
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
			return
		}

		log.Info("\n\n\nOrder body", string(body))

		// ✅ Log file path
		logFile := "shopify_orders.log"

		// ✅ Create or append to the log file
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open log file"})
			return
		}
		defer file.Close()

		// ✅ Write timestamp and order data
		logEntry := fmt.Sprintf("\n\n=== Shopify Order - %s ===\n%s\n", time.Now().Format(time.RFC3339), string(body))
		if _, err := file.WriteString(logEntry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write log"})
			return
		}

		// ✅ Print to console (for debugging)
		fmt.Println("🛍️ Shopify Order Created:", string(body))

		// ✅ Respond to Shopify (must be 200 within 5 seconds)
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
