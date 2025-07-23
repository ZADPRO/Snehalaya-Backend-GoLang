package purchaseOrderRoutes

import (
	purchaseOrderController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"

)

func PurhcaseOrderRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/purchaseOrder")

	// CREATE INITIAL PRODUCTS
	route.POST("/create", accesstoken.JWTMiddleware(), purchaseOrderController.CreatePurchaseOrderController())
	route.GET("/read", accesstoken.JWTMiddleware(), purchaseOrderController.GetAllPurchaseOrdersController())
	route.GET("/read/:id", accesstoken.JWTMiddleware(), purchaseOrderController.GetPurchaseOrderByIdController())

	route.GET("/dummy-products/:purchaseOrderId", accesstoken.JWTMiddleware(), purchaseOrderController.GetDummyProductsByPOID())
	// UPDATE PURCHASE ORDER PRODUCTS
	route.PUT("/dummy-products/update", accesstoken.JWTMiddleware(), purchaseOrderController.UpdateDummyProductStatus())
	// BULK UPDATE - ACCEPT, REJECT, UNDO
	route.PUT("/dummy-products/bulk-accept", accesstoken.JWTMiddleware(), purchaseOrderController.BulkAcceptDummyProducts())
	route.PUT("/dummy-products/bulk-reject", accesstoken.JWTMiddleware(), purchaseOrderController.BulkRejectDummyProducts())
	route.PUT("/dummy-products/bulk-undo", accesstoken.JWTMiddleware(), purchaseOrderController.BulkUndoDummyProducts())

}
