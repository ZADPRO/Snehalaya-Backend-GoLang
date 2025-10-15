package oldProductRoutes

import (
	oldProductController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/oldProductMigration/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

func OldProductMigrationRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/oldProductMigration")
	route.POST("/create", accesstoken.JWTMiddleware(), oldProductController.MigrateOldProductsController())
}
