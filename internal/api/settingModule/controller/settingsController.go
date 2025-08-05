package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/model"
	settingsService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	roleType "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/GetRoleType"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"
)

// CATEGORIES CONTROLLER

func CreateCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Create Category Controller")

		// Fetching user-related data from context
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Info("\n\nRole ID Console", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var category model.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			log.Error("Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		// Fetch role name from DB
		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		fmt.Println("roleName", roleName)
		if err != nil {
			log.Error("Failed to get role name: " + err.Error())
		} else {
			log.Info("Role Name: " + roleName)
		}

		// Create category
		err = settingsService.CreateCategoryService(dbConnt, &category, roleName)
		if err != nil {
			log.Error("Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create category"})
			}
			return
		}

		// Success
		log.Info("Category created successfully")
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Category created successfully",
			"token":   token,
		})
	}
}

func GetAllCategoriesController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Get All Categories Controller")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		categories := settingsService.GetAllCategoriesService(dbConnt)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "data": categories, "token": token})
	}
}

func UpdateCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Update Category Controller")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		var category model.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			log.Error("Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		// Fetch role name from DB
		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		fmt.Println("roleName", roleName)
		if err != nil {
			log.Error("Failed to get role name: " + err.Error())
		} else {
			log.Info("Role Name: " + roleName)
		}

		errH := settingsService.UpdateCategoryService(dbConnt, &category, roleName)
		if errH != nil {
			log.Error("Service error: " + errH.Error())
			if errH.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update category"})
			}
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("Category updated successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category updated successfully", "token": token})
	}
}

func DeleteCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Delete Category Controller")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		categoryId := c.Param("id")
		forceDelete := c.DefaultQuery("forceDelete", "false") == "true"

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		subcategories, err := settingsService.GetSubcategoriesByCategory(dbConnt, categoryId)
		if err != nil {
			log.Error("Error fetching subcategories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Internal server error"})
			return
		}

		if len(subcategories) > 0 && !forceDelete {
			// Return subcategories and ask for confirmation
			c.JSON(http.StatusConflict, gin.H{
				"status":             false,
				"message":            "This category contains subcategories. Deleting it will make them idle.",
				"subcategories":      subcategories,
				"confirmationNeeded": true,
			})
			return
		}

		err = settingsService.DeleteCategoryService(dbConnt, categoryId)
		if err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete category"})
			return
		}

		log.Info("Category deleted successfully")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category deleted successfully", "token": token})
	}
}

func BulkDeleteCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Bulk Delete Category Controller")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var request struct {
			CategoryIDs []int `json:"categoryIds"`
			ForceDelete bool  `json:"forceDelete"`
		}

		if err := c.ShouldBindJSON(&request); err != nil || len(request.CategoryIDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid category IDs"})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// ✅ Use CheckSubcategoriesExistence here
		subcategoriesMap, err := settingsService.CheckSubcategoriesExistence(dbConnt, request.CategoryIDs)
		if err != nil {
			log.Error("Error checking subcategories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Internal server error"})
			return
		}

		if len(subcategoriesMap) > 0 && !request.ForceDelete {
			c.JSON(http.StatusConflict, gin.H{
				"status":             false,
				"message":            "Some categories contain subcategories. Deleting them will make subcategories idle.",
				"subcategoriesMap":   subcategoriesMap,
				"confirmationNeeded": true,
			})
			return
		}

		err = settingsService.BulkDeleteCategoriesService(dbConnt, request.CategoryIDs)
		if err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete categories"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Categories deleted successfully",
			"token":   token,
		})
	}
}

// SUB CATEGORIES CONTROLLER

func CreateSubCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Create SubCategory Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		var subCateogry model.SubCategory
		if err := c.ShouldBindJSON(&subCateogry); err != nil {
			log.Error("Invalid request body : " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		if err := settingsService.CreateSubCategoryService(dbConnt, &subCateogry); err != nil {
			log.Error("Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create sub category"})
			}
			return
		}

		log.Info("Sub category created successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Sub category created", "token": token})
	}
}

func GetAllSubCategoriesController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Get All SubCategories Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data := settingsService.GetAllSubCategoriesService(dbConnt)
		log.Info("Fetched subcategories: ", data)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "data": data, "token": token})
	}
}

func UpdateSubCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Update SubCategory Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		var sub model.SubCategory
		if err := c.ShouldBindJSON(&sub); err != nil {
			log.Error("Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := settingsService.UpdateSubCategoryService(dbConnt, &sub); err != nil {
			log.Error("Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update sub category"})
			}
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("Sub category updated successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Sub category updated", "token": token})
	}
}

func DeleteSubCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Delete SubCategory Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		id := c.Param("id")
		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := settingsService.DeleteSubCategoryService(dbConnt, id); err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete sub category"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("Sub category deleted successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Sub category deleted", "token": token})
	}
}

// BRANCHES CONTROLLER
func CreateBranchController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Create Branch Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var branch model.Branch
		if err := c.ShouldBindJSON(&branch); err != nil {
			log.Error("Invalid request body : " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		if err := settingsService.CreateBranchService(dbConnt, &branch); err != nil {
			log.Error("Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create branch"})
			}
			return
		}

		log.Info("Branch created successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch created", "token": token})
	}
}

func GetAllBranchesController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Get All Branches Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found in request context."})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		branches := settingsService.GetAllBranchesService(dbConnt)
		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{"status": true, "data": branches, "token": token})
	}
}

func UpdateBranchController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Update Branch Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found in request context."})
			return
		}

		var branch model.Branch
		if err := c.ShouldBindJSON(&branch); err != nil {
			log.Error("Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := settingsService.UpdateBranchService(dbConnt, &branch); err != nil {
			log.Error("Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update branch"})
			}
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch updated", "token": token})
	}
}

func DeleteBranchController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Delete Branch Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found in request context."})
			return
		}

		id := c.Param("id")
		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := settingsService.DeleteBranchService(dbConnt, id); err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete branch"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch deleted", "token": token})
	}
}

// BRANCH WITH FLOOR CONTROLLER
func CreateNewBranchWithFloorController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("\n\nCreate Branch Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var payload struct {
			model.BranchWithFloor
			Floors []struct {
				FloorName string
				FloorCode string
				Sections  []struct {
					CategoryId       int
					RefSubCategoryId int
					SectionName      string
					SectionCode      string
				}
			}
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Error("Invalid payload: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request payload"})
			return
		}

		err := settingsService.CreateNewBranchWithFloor(dbConnt, &payload.BranchWithFloor, payload.Floors, idValue.(int))
		if err != nil {
			log.Error("Failed to create branch with floors: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		log.Info("Branch with Floors and Sections created successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch created with Floors and Sections", "token": token})
	}
}

// ATTRIBUTES
// func CreateAttributeGroupController() gin.HandlerFunc {

// }

// ADD NEW EMPLOYEE CONTROLLER
func GetEmployeeRoleType() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Get Employee Role Type Controller ===> ")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found in request context."})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleTypes := settingsService.GetUserRoleTypeService(dbConn)

		if roleTypes == nil || len(roleTypes) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "No role types found"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Role types fetched successfully",
			"roles":   roleTypes,
			"token":   token,
		})
	}
}

func CreateEmployeeController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Create Employee Controller")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")
		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Missing context info"})
			return
		}

		var payload model.EmployeePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Error("Invalid JSON: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		err := settingsService.CreateEmployeeService(dbConn, &payload)
		if err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		log.Info("Employee created successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Employee created", "token": token})
	}
}

func GetAllEmployeesController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Create Employee Controller")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")
		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Missing context info"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		employees, err := settingsService.GetAllEmployeesService(dbConn)
		if err != nil {
			log.Error("Failed to fetch employees: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}
		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		log.Info("Fetched all employees successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "data": employees, "token": token})
	}
}

func GetEmployeeByIDController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		id := c.Param("id")
		employee, err := settingsService.GetEmployeeByIDService(dbConn, id)
		if err != nil {
			log.Error("Failed to fetch employee: " + err.Error())
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "data": employee})
	}
}

func UpdateEmployeeController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		id := c.Param("id")
		var payload model.EmployeePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		err := settingsService.UpdateEmployeeService(dbConn, id, &payload)
		if err != nil {
			log.Error("Failed to update employee: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Employee updated successfully"})
	}
}

func DeleteEmployeeController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		id := c.Param("id")
		err := settingsService.SoftDeleteEmployeeService(dbConn, id)
		if err != nil {
			log.Error("Failed to delete employee: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Employee deleted (soft) successfully"})
	}
}

func GetEmployeeController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		idValue, exists := c.Get("id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized: No ID in token"})
			return
		}

		idStr := fmt.Sprintf("%v", idValue) // convert to string if needed

		employee, err := settingsService.GetEmployeeService(dbConn, idStr)
		if err != nil {
			log.Error("Failed to fetch employee: " + err.Error())
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "data": employee})
	}
}

func UpdateProfileController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// ✅ Get user ID from token (context set by middleware)
		idValue, exists := c.Get("id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
			return
		}

		var id string
		switch v := idValue.(type) {
		case string:
			id = v
		case float64:
			id = fmt.Sprintf("%.0f", v)
		case int:
			id = strconv.Itoa(v)
		default:
			log.Error(fmt.Sprintf("Unexpected ID type: %T", v))
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Invalid user ID format"})
			return
		}

		var payload model.ProfilePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		err := settingsService.UpdateProfileService(dbConn, id, &payload)
		if err != nil {
			log.Error("Failed to update Profile: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Profile updated successfully"})
	}
}
