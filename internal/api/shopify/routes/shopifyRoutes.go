package shopfiyRoutes

import (
	shopifyController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/shopify/controller"
	"github.com/gin-gonic/gin"
)

func RegisterShopifyRoutes(router *gin.Engine) {
	api := router.Group("/api/v1/shopify")
	{
		api.GET("/products", shopifyController.GetShopifyProducts)
		api.POST("/products", shopifyController.CreateShopifyProduct)
	}

	api2 := router.Group("/api/v1/webhook/shopify")
	{
		api2.POST("/ordercreation", shopifyController.OrderCreationWebhook())
	}
}
