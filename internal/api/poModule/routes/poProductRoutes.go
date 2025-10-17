package PORoutes

import (
	poController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

func PurchaseOrderProductRoutes(route *gin.Engine) {
	poGroup := route.Group("/api/v1/admin")
	{
		poGroup.POST("/poProductsUpdate", accesstoken.JWTMiddleware(), poController.NewPurchaseOrderController().CreatePurchaseOrderProductsController())
	}
}
