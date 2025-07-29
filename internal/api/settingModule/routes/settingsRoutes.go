package settingsRoutes

import (
	settingsController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

func SettingsAdminRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/settings")

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

	// BRANCHES ROUTES
	route.POST("/branches", accesstoken.JWTMiddleware(), settingsController.CreateBranchController())
	route.GET("/branches", accesstoken.JWTMiddleware(), settingsController.GetAllBranchesController())
	route.PUT("/branches", accesstoken.JWTMiddleware(), settingsController.UpdateBranchController())
	route.DELETE("/branches/:id", accesstoken.JWTMiddleware(), settingsController.DeleteBranchController())

	// USER ROLES

	// ATTRIBUTES
	// route.POST("/attributes", accesstoken.JWTMiddleware(), settingsController.CreateAttributeGroupController())

	// EMPLOYEES ROUTES
	route.GET("/employeeRoleType", accesstoken.JWTMiddleware(), settingsController.GetEmployeeRoleType())
	route.POST("/employees", accesstoken.JWTMiddleware(), settingsController.CreateEmployeeController())
	route.GET("/employees", accesstoken.JWTMiddleware(), settingsController.GetAllEmployeesController())
	route.GET("/employees/:id", accesstoken.JWTMiddleware(), settingsController.GetEmployeeByIDController())
	route.PUT("/employees/:id", accesstoken.JWTMiddleware(), settingsController.UpdateEmployeeController())
	route.DELETE("/employees/:id", accesstoken.JWTMiddleware(), settingsController.DeleteEmployeeController())
	route.GET("/getEmployees", accesstoken.JWTMiddleware(), settingsController.GetEmployeeController())
	route.GET("/updateEmployeeProfile", accesstoken.JWTMiddleware(), settingsController.UpdateProfileController())
	
}
