package shopifyController

import (
	shopifyHelper "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/shopify/helper"
	shopifyService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/shopify/service"
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
