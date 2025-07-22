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

}
