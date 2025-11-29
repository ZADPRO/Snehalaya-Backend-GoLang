package settingsRoutes

import (
	settingsController "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/controller"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"

)

func SettingsAdminRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/settings")
	routev2 := router.Group("/api/v2/admin/settings")

	route.POST("/initialCategoryCode", accesstoken.JWTMiddleware(), settingsController.CheckInitialCategoryCodeController())

	// INITIAL ROUTES
	route.POST("/initialCategories", accesstoken.JWTMiddleware(), settingsController.CreateInitialCategoryController())
	route.GET("/initialCategories", accesstoken.JWTMiddleware(), settingsController.GetAllInitialCategoryController())
	route.PUT("/initialCategories", accesstoken.JWTMiddleware(), settingsController.UpdateInitialCategoryController())
	route.DELETE("/initialCategories", accesstoken.JWTMiddleware(), settingsController.DeleteInitialCategoryController())

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
	routev2.GET("/branches/:id", accesstoken.JWTMiddleware(), settingsController.GetNewBranchWithFloorWithIdController())
	routev2.PUT("/branches/:id", accesstoken.JWTMiddleware(), settingsController.UpdateBranchWithFloorController())
	routev2.DELETE("/branches/:id", accesstoken.JWTMiddleware(), settingsController.SoftDeleteBranchController())

	// USER ROLES

	// ATTRIBUTES
	// route.POST("/attributes", accesstoken.JWTMiddleware(), settingsController.CreateAttributeGroupController())
	route.GET("/attributesDataType", accesstoken.JWTMiddleware(), settingsController.GetAttributeDataType())
	route.POST("/attributes", accesstoken.JWTMiddleware(), settingsController.CreateAttributeGroupController())
	route.GET("/attributes", accesstoken.JWTMiddleware(), settingsController.GetAttributeGroupController())
	route.PUT("/attributes", accesstoken.JWTMiddleware(), settingsController.UpdateAttributeGroupController())
	route.POST("/attributesHide", accesstoken.JWTMiddleware(), settingsController.DeleteAttributeGroupController())

	routev2.POST("/attributes")

	route.POST("/product-fields", accesstoken.JWTMiddleware(), settingsController.CreateProductFieldController())
	route.GET("/product-fields", accesstoken.JWTMiddleware(), settingsController.GetAllProductFieldsController())
	route.PUT("/product-fields", accesstoken.JWTMiddleware(), settingsController.UpdateProductFieldController())

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

	// SETTINGS PRODUCTS ROUTES
	route.POST("/settingsProducts", accesstoken.JWTMiddleware(), settingsController.CreateSettingsProductController())
	route.GET("/settingsProducts", accesstoken.JWTMiddleware(), settingsController.GetAllSettingsProductsController())
	route.PUT("/settingsProducts", accesstoken.JWTMiddleware(), settingsController.UpdateSettingsProductController())
	route.DELETE("/settingsProducts", accesstoken.JWTMiddleware(), settingsController.DeleteSettingsProductsController())

	// DESIGN
	route.POST("/design", accesstoken.JWTMiddleware(), settingsController.CreateMasterController("design"))
	route.GET("/design", accesstoken.JWTMiddleware(), settingsController.GetAllMasterController("design"))
	route.PUT("/design", accesstoken.JWTMiddleware(), settingsController.UpdateMasterController("design"))
	route.DELETE("/design", accesstoken.JWTMiddleware(), settingsController.DeleteMasterController("design"))

	// COLOR
	route.POST("/color", accesstoken.JWTMiddleware(), settingsController.CreateMasterController("color"))
	route.GET("/color", accesstoken.JWTMiddleware(), settingsController.GetAllMasterController("color"))
	route.PUT("/color", accesstoken.JWTMiddleware(), settingsController.UpdateMasterController("color"))
	route.DELETE("/color", accesstoken.JWTMiddleware(), settingsController.DeleteMasterController("color"))

	// BRAND
	route.POST("/brand", accesstoken.JWTMiddleware(), settingsController.CreateMasterController("brand"))
	route.GET("/brand", accesstoken.JWTMiddleware(), settingsController.GetAllMasterController("brand"))
	route.PUT("/brand", accesstoken.JWTMiddleware(), settingsController.UpdateMasterController("brand"))
	route.DELETE("/brand", accesstoken.JWTMiddleware(), settingsController.DeleteMasterController("brand"))

	// SIZE
	route.POST("/size", accesstoken.JWTMiddleware(), settingsController.CreateMasterController("size"))
	route.GET("/size", accesstoken.JWTMiddleware(), settingsController.GetAllMasterController("size"))
	route.PUT("/size", accesstoken.JWTMiddleware(), settingsController.UpdateMasterController("size"))
	route.DELETE("/size", accesstoken.JWTMiddleware(), settingsController.DeleteMasterController("size"))

	// VARIENT
	route.POST("/varient", accesstoken.JWTMiddleware(), settingsController.CreateMasterController("Varient"))
	route.GET("/varient", accesstoken.JWTMiddleware(), settingsController.GetAllMasterController("Varient"))
	route.PUT("/varient", accesstoken.JWTMiddleware(), settingsController.UpdateMasterController("Varient"))
	route.DELETE("/varient", accesstoken.JWTMiddleware(), settingsController.DeleteMasterController("Varient"))

	// PATTERNS
	route.POST("/patterns", accesstoken.JWTMiddleware(), settingsController.CreateMasterController("Patterns"))
	route.GET("/patterns", accesstoken.JWTMiddleware(), settingsController.GetAllMasterController("Patterns"))
	route.PUT("/patterns", accesstoken.JWTMiddleware(), settingsController.UpdateMasterController("Patterns"))
	route.DELETE("/patterns", accesstoken.JWTMiddleware(), settingsController.DeleteMasterController("Patterns"))

}
