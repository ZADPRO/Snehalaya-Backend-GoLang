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

	// PURCHASE ORDER - VIEW ALL PRODUCTS
	route.GET("/list-all-products", accesstoken.JWTMiddleware(), purchaseOrderController.GetReceivedDummyProductsController())
	route.GET("/list-all-products-barcode", accesstoken.JWTMiddleware(), purchaseOrderController.GetReceivedDummyProductsBarcodeController())

	// CREATE CATALOG
	route.POST("/products", accesstoken.JWTMiddleware(), purchaseOrderController.CreateProductController())

	// LATEST CHANGES FOR PO CREATION
	route.POST("/createPurchaseOrder", accesstoken.JWTMiddleware(), purchaseOrderController.NewCreatePurchaseOrderController())
	route.GET("/getOurchaseOrder", accesstoken.JWTMiddleware(), purchaseOrderController.NewGetAllPurchaseOrdersController())
	route.GET("/purchaseOrder/:id", accesstoken.JWTMiddleware(), purchaseOrderController.NewGetSinglePurchaseOrderController())

	// GRN
	route.POST("/createGRN", accesstoken.JWTMiddleware(), purchaseOrderController.NewCreateGRNController())
	route.GET("/grn/list", accesstoken.JWTMiddleware(), purchaseOrderController.NewGetAllGRNController())
	route.GET("/grn/:id", accesstoken.JWTMiddleware(), purchaseOrderController.NewGetSingleGRNController())

	// INVENTORY
	route.GET("/getInventoryList",
		accesstoken.JWTMiddleware(),
		purchaseOrderController.NewGetInventoryListController(),
	)
}
