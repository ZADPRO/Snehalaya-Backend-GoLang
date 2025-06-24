package settingsRoutes

import (
	settingsController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"

)

func SettingsAdminRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/settings")

	route.POST("/categories", accesstoken.JWTMiddleware(), settingsController.CreateCategoryController())
	route.GET("/categories", accesstoken.JWTMiddleware(), settingsController.GetAllCategoriesController())
	route.PUT("/categories", accesstoken.JWTMiddleware(), settingsController.UpdateCategoryController())
	route.DELETE("/categories/:id", accesstoken.JWTMiddleware(), settingsController.DeleteCategoryController())
}
