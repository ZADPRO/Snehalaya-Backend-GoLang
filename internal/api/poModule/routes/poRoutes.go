package PORoutes

import (
	poController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

func PurchaseOrderRoutes(route *gin.Engine) {
	po := route.Group("/api/v1/admin")
	{
		po.POST("/purchaseOrder", accesstoken.JWTMiddleware(), poController.CreatePurchaseOrderController())
		po.GET("/purchaseOrder", accesstoken.JWTMiddleware(), poController.GetAllPurchaseOrdersController())
		po.PUT("/purchaseOrder", accesstoken.JWTMiddleware(), poController.UpdatePurchaseOrderController())
		// po.DELETE("/:id", accesstoken.JWTMiddleware(), poController.DeletePurchaseOrderController())

		po.GET("/getAllPurchaseOrders", accesstoken.JWTMiddleware(), poController.GetAllPurchaseOrdersListController())

		po.POST("/updatePurchaseOrderProducts", accesstoken.JWTMiddleware(), poController.UpdatePurchaseOrderProductsController())

		po.POST("/savePurchaseOrderProducts", accesstoken.JWTMiddleware(), poController.SavePurchaseOrderProductsController())

		po.GET("/getAcceptedProducts/:purchaseOrderId", accesstoken.JWTMiddleware(), poController.GetAcceptedProductsController())

		po.GET("/details/:purchaseOrderNumber", accesstoken.JWTMiddleware(), poController.GetPurchaseOrderDetailsController())

	}
}
