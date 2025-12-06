package productRoutes

import (
	productController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/products/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"

)

func ProductManagementRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/products")
	route.POST("/create", accesstoken.JWTMiddleware(), productController.CreatePOProductController())
	route.GET("/read", accesstoken.JWTMiddleware(), productController.GetAllPOProductsController())
	route.GET("/read/:id", accesstoken.JWTMiddleware(), productController.GetPOProductByIdController())
	route.PUT("/update", accesstoken.JWTMiddleware(), productController.UpdatePOProductController())
	route.DELETE("/delete/:id", accesstoken.JWTMiddleware(), productController.DeletePOProductController())
	route.POST("/check-sku", accesstoken.JWTMiddleware(), productController.CheckSKUInBranchController())
	route.GET("/branch-4-products", accesstoken.JWTMiddleware(), productController.GetBranch4ProductsController())

	// INVENTORY STOCK TRANSFER
	route.POST("/stock-transfer", accesstoken.JWTMiddleware(), productController.CreateStockTransfer())
	route.GET("/stock-transfer", accesstoken.JWTMiddleware(), productController.GetStockTransfersController())

	route.GET("/stock-transfer/all", accesstoken.JWTMiddleware(), productController.GetAllStockTransfersController())

	route.PUT("/stock-transfer/receive", accesstoken.JWTMiddleware(), productController.ReceiveStockProductsController())

	route.POST("/save", accesstoken.JWTMiddleware(), productController.SaveProductImagesController())
	route.GET("/byProduct/:productInstanceId", productController.GetImagesByProductController())

	route.GET("/purchaseOrderAcceptedProducts/:id", accesstoken.JWTMiddleware(), productController.GetSinglePurchaseOrderAcceptedProductController())

	route.POST(
		"/check-sku-grn",
		accesstoken.JWTMiddleware(),
		productController.CheckSKUInGRNController(),
	)

	route.POST("/new-stock-transfer", accesstoken.JWTMiddleware(), productController.StockTransferController())

	route.GET("/stock-transfer/list", accesstoken.JWTMiddleware(), productController.GetStockTransferMasterController())

	route.GET("/stock-transfer/items/:transferId", accesstoken.JWTMiddleware(), productController.GetStockTransferItemsController())

	// route.GET("/stock-transfer/:id", accesstoken.JWTMiddleware(), productController.GetStockTransferByIDController())

	route.POST("/createBundle", accesstoken.JWTMiddleware(), productController.CreateBundleInwardController())
	route.GET("/getBundle", accesstoken.JWTMiddleware(), productController.GetAllBundleInwardsController())
	route.PUT("/updateBundle", accesstoken.JWTMiddleware(), productController.UpdateBundleInwardController())

}
