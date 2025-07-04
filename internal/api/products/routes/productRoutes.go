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
}
