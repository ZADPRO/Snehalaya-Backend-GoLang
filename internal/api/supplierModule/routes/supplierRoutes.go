package supplierRoutes

import (
	supplierController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"

)

func SupplierRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/suppliers")

	route.POST("/create", accesstoken.JWTMiddleware(), supplierController.CreateSupplierController())
	route.GET("/read", accesstoken.JWTMiddleware(), supplierController.GetAllSuppliersController())
	route.GET("/read/:id", accesstoken.JWTMiddleware(), supplierController.GetSupplierByIdController())
	route.PUT("/update", accesstoken.JWTMiddleware(), supplierController.UpdateSupplierController())
	route.DELETE("/delete/:id", accesstoken.JWTMiddleware(), supplierController.DeleteSupplierController())
	route.DELETE("/delete", accesstoken.JWTMiddleware(), supplierController.DeleteSupplierController())

	route.POST("/delete/bulk", accesstoken.JWTMiddleware(), supplierController.BulkDeleteSupplierController())

}
