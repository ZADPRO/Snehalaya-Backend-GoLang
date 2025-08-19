package settingsRoutes

import (
	settingsController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"

)

func SettingsAdminRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/settings")
	routev2 := router.Group("/api/v2/admin/settings")
	// CATEGORIES ROUTES
	route.POST("/categories", accesstoken.JWTMiddleware(), settingsController.CreateCategoryController())
	route.GET("/categories", accesstoken.JWTMiddleware(), settingsController.GetAllCategoriesController())
	route.PUT("/categories", accesstoken.JWTMiddleware(), settingsController.UpdateCategoryController())
	route.DELETE("/categories/:id", accesstoken.JWTMiddleware(), settingsController.DeleteCategoryController())
	route.DELETE("/categories", accesstoken.JWTMiddleware(), settingsController.BulkDeleteCategoryController())

	// SUB CATEGORIES ROUTES
	route.POST("/subcategories", accesstoken.JWTMiddleware(), settingsController.CreateSubCategoryController())
	route.GET("/subcategories", accesstoken.JWTMiddleware(), settingsController.GetAllSubCategoriesController())
	route.PUT("/subcategories", accesstoken.JWTMiddleware(), settingsController.UpdateSubCategoryController())
	route.DELETE("/subcategories/:id", accesstoken.JWTMiddleware(), settingsController.DeleteSubCategoryController())
	route.DELETE("/subcategories", accesstoken.JWTMiddleware(), settingsController.BulkDeleteSubCategoryController())

	// BRANCHES ROUTES
	route.POST("/branches", accesstoken.JWTMiddleware(), settingsController.CreateBranchController())
	route.GET("/branches", accesstoken.JWTMiddleware(), settingsController.GetAllBranchesController())
	route.PUT("/branches", accesstoken.JWTMiddleware(), settingsController.UpdateBranchController())
	route.DELETE("/branches/:id", accesstoken.JWTMiddleware(), settingsController.DeleteBranchController())

	// BRANCH WITH FLOOR ROUTES
	routev2.POST("/branches", accesstoken.JWTMiddleware(), settingsController.CreateNewBranchWithFloorController())
	routev2.GET("/branches", accesstoken.JWTMiddleware(), settingsController.GetNewBranchWithFloorController())

	// USER ROLES

	// ATTRIBUTES
	// route.POST("/attributes", accesstoken.JWTMiddleware(), settingsController.CreateAttributeGroupController())
	route.GET("/attributesDataType", accesstoken.JWTMiddleware(), settingsController.GetAttributeDataType())
	route.POST("/attributes", accesstoken.JWTMiddleware(), settingsController.CreateAttributeGroupController())
	route.GET("/attributes", accesstoken.JWTMiddleware(), settingsController.GetAttributeGroupController())
	route.PUT("/attributes", accesstoken.JWTMiddleware(), settingsController.UpdateAttributeGroupController())
	route.POST("/attributesHide", accesstoken.JWTMiddleware(), settingsController.DeleteAttributeGroupController())

	// EMPLOYEES ROUTES
	route.GET("/employeeRoleType", accesstoken.JWTMiddleware(), settingsController.GetEmployeeRoleType())
	route.POST("/employees", accesstoken.JWTMiddleware(), settingsController.CreateEmployeeController())
	route.GET("/employees", accesstoken.JWTMiddleware(), settingsController.GetAllEmployeesController())
	route.GET("/employees/:id", accesstoken.JWTMiddleware(), settingsController.GetEmployeeByIDController())
	route.PUT("/employees/:id", accesstoken.JWTMiddleware(), settingsController.UpdateEmployeeController())
	route.DELETE("/employees/:id", accesstoken.JWTMiddleware(), settingsController.DeleteEmployeeController())
	route.GET("/getEmployees", accesstoken.JWTMiddleware(), settingsController.GetEmployeeController())
	route.PUT("/updateEmployeeProfile", accesstoken.JWTMiddleware(), settingsController.UpdateProfileController())

	// overView
	route.GET("/overview", accesstoken.JWTMiddleware(), settingsController.GetSettingsOverview())

}
