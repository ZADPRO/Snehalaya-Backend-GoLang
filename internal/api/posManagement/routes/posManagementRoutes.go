package  posManagementRoutes

import (
	posManagementController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/posManagement/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)


func POSManagementRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/pos")
	route.POST("/customer", accesstoken.JWTMiddleware(), posManagementController.AddCustomer())
	// route.GET("/read", accesstoken.JWTMiddleware(), productController.GetAllPOProductsController())
	// route.GET("/read/:id", accesstoken.JWTMiddleware(), productController.GetPOProductByIdController())
	// route.PUT("/update", accesstoken.JWTMiddleware(), productController.UpdatePOProductController())
	// route.DELETE("/delete/:id", accesstoken.JWTMiddleware(), productController.DeletePOProductController())
}
