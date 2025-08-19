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
		log.Info("\n\n\n🚀 Create Category Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var category model.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Request Body: %+v", category)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.CreateCategoryService(dbConnt, &category, roleName)
		if err != nil {
			log.Error("❌ Service Error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create category"})
			}
			return
		}

		log.Info("✅ Category created successfully\n\n")
		log.Info("\n=================================================================\n")

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
		log.Info("\n\n📥 GetAllCategoriesController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			err := sqlDB.Close()
			if err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		log.Info("📦 Fetching all categories from DB")
		categories := settingsService.GetAllCategoriesService(dbConnt)
		log.Infof("📊 Categories fetched: count = %d", len(categories))

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Sending response with category list\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   categories,
			"token":  token,
		})
	}
}

func UpdateCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 UpdateCategoryController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var category model.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			log.Error("❌ Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Request Body: %+v", category)

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("❌ Failed to get role name: " + err.Error())
		} else {
			log.Infof("👤 Role Name: %s", roleName)
		}

		log.Info("🛠️ Calling UpdateCategoryService")
		errH := settingsService.UpdateCategoryService(dbConnt, &category, roleName)
		if errH != nil {
			log.Error("❌ Service error: " + errH.Error())
			if errH.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update category"})
			}
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Category updated successfully\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Category updated successfully",
			"token":   token,
		})
	}
}

func DeleteCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 DeleteCategoryController invoked")

		// Get context values
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing user/role/branch information in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		categoryId := c.Param("id")
		forceDelete := c.DefaultQuery("forceDelete", "false") == "true"
		log.Infof("🗑️ Delete Request: categoryId=%s, forceDelete=%t", categoryId, forceDelete)

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		// Check subcategories
		log.Info("🔎 Checking for subcategories before deletion")
		subcategories, err := settingsService.GetSubcategoriesByCategory(dbConnt, categoryId)
		if err != nil {
			log.Error("❌ Error fetching subcategories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Internal server error",
			})
			return
		}

		if len(subcategories) > 0 && !forceDelete {
			log.Warn("⚠️ Subcategories found. Confirmation required before force delete.")
			c.JSON(http.StatusConflict, gin.H{
				"status":             false,
				"message":            "This category contains subcategories. Deleting it will make them idle.",
				"subcategories":      subcategories,
				"confirmationNeeded": true,
			})
			return
		}

		// Perform deletion
		log.Info("🛠️ Calling DeleteCategoryService")
		err = settingsService.DeleteCategoryService(dbConnt, categoryId)
		if err != nil {
			log.Error("❌ Service error during category deletion: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to delete category",
			})
			return
		}

		log.Info("✅ Category deleted successfully\n\n")
		log.Info("\n=================================================================\n")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Category deleted successfully",
			"token":   token,
		})
	}
}

func BulkDeleteCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 BulkDeleteCategoryController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
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
			log.Error("❌ Invalid request body or empty category IDs")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid category IDs",
			})
			return
		}
		log.Infof("📦 Bulk delete request: categoryIds=%v, forceDelete=%v", request.CategoryIDs, request.ForceDelete)

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		// Step 1: Check for subcategories
		log.Info("🔎 Checking for subcategories in selected categories")
		subcategoriesMap, err := settingsService.CheckSubcategoriesExistence(dbConnt, request.CategoryIDs)
		if err != nil {
			log.Error("❌ Error checking subcategories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Internal server error",
			})
			return
		}

		if len(subcategoriesMap) > 0 && !request.ForceDelete {
			log.Warn("⚠️ Some categories have subcategories. Confirmation needed before force delete.")
			c.JSON(http.StatusConflict, gin.H{
				"status":             false,
				"message":            "Some categories contain subcategories. Deleting them will make subcategories idle.",
				"subcategoriesMap":   subcategoriesMap,
				"confirmationNeeded": true,
			})
			return
		}

		// Step 2: Perform soft delete
		log.Info("🛠️ Proceeding to soft delete categories")

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.BulkDeleteCategoriesService(dbConnt, request.CategoryIDs, roleName)
		if err != nil {
			log.Error("❌ Service error during bulk delete: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to delete categories",
			})
			return
		}

		// Optional: you can loop and log each deletion
		log.Infof("✅ Categories soft deleted successfully: %v\n\n", request.CategoryIDs)
		log.Info("\n=================================================================\n")

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
		log.Info("\n\n\n🚀 Create SubCategory Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var subCategory model.SubCategory
		if err := c.ShouldBindJSON(&subCategory); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Request Body: %+v", subCategory)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		// ✅ Extract role name
		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		if err := settingsService.CreateSubCategoryService(dbConnt, &subCategory, roleName); err != nil {
			log.Error("❌ Service Error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create sub category"})
			}
			return
		}

		log.Info("✅ SubCategory created successfully")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Sub category created successfully",
			"token":   token,
		})
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
		log.Info("\n\n\n🚀 Update SubCategory Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var sub model.SubCategory
		if err := c.ShouldBindJSON(&sub); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📥 Input SubCategory: %+v", sub)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
			roleName = "Unknown"
		} else {
			log.Infof("👤 Role Name resolved: %s", roleName)
		}

		if err := settingsService.UpdateSubCategoryService(dbConnt, &sub, roleName); err != nil {
			log.Error("❌ Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update sub category"})
			}
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ SubCategory updated successfully")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Sub category updated",
			"token":   token,
		})
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

func BulkDeleteSubCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 BulkDeleteSubCategoryController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var request struct {
			SubCategoryIDs []int `json:"subCategoriesId"`
		}

		if err := c.ShouldBindJSON(&request); err != nil || len(request.SubCategoryIDs) == 0 {
			log.Error("❌ Invalid request body or empty subcategory IDs")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid subcategory IDs",
			})
			return
		}

		log.Infof("📦 Bulk delete request: subCategoryIds=%v", request.SubCategoryIDs)

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.BulkDeleteSubCategoriesService(dbConnt, request.SubCategoryIDs, roleName)
		if err != nil {
			log.Error("❌ Service error during bulk subcategory delete: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to delete subcategories",
			})
			return
		}

		log.Infof("✅ SubCategories soft deleted successfully: %v\n\n", request.SubCategoryIDs)
		log.Info("\n=================================================================\n")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Subcategories deleted successfully",
			"token":   token,
		})
	}
}

// BRANCHES CONTROLLER
func CreateBranchController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n\n🚀 Create Branch Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var branch model.Branch
		if err := c.ShouldBindJSON(&branch); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Request Body: %+v", branch)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.CreateBranchService(dbConnt, &branch, roleName)
		if err != nil {
			log.Error("❌ Service Error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create branch"})
			}
			return
		}

		log.Info("✅ Branch created successfully\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Branch created successfully",
			"token":   token,
		})
	}
}

func GetAllBranchesController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n\n📥 Get All Branches Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		branches, err := settingsService.GetAllBranchesService(dbConnt)
		if err != nil {
			log.Error("❌ Failed to get branches: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch branches",
			})
			return
		}

		log.Infof("✅ %d branches retrieved", len(branches))

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   branches,
			"token":  token,
		})
	}
}

func UpdateBranchController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n\n🔧 Update Branch Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var branch model.Branch
		if err := c.ShouldBindJSON(&branch); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		log.Infof("📥 Request Body: %+v", branch)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.UpdateBranchService(dbConnt, &branch, roleName)
		if err != nil {
			log.Error("❌ Service error: " + err.Error())

			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update branch"})
			}
			return
		}

		log.Info("✅ Branch updated successfully")
		log.Info("\n=================================================================\n")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch updated successfully", "token": token})
	}
}

func DeleteBranchController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n\n🗑️ Delete Branch Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values")
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found in request context."})
			return
		}

		id := c.Param("id")
		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
			roleName = "Unknown"
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.DeleteBranchService(dbConnt, id, roleName)
		if err != nil {
			log.Error("❌ Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete branch"})
			return
		}

		log.Info("✅ Branch soft deleted successfully")
		log.Info("\n=================================================================\n")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch deleted successfully", "token": token})
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

		userId := 0
		switch v := idValue.(type) {
		case float64:
			userId = int(v)
		case int:
			userId = v
		default:
			// handle unexpected type case
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Invalid user ID type"})
			return
		}

		// ✅ Extract role name
		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, errRole := roleType.GetRoleTypeNameByID(dbConnt, roleId)

		if errRole != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.CreateNewBranchWithFloor(dbConnt, &payload.BranchWithFloor, payload.Floors, userId)
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

func GetNewBranchWithFloorController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("\n\n📥 GetNewBranchWithFloorController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v (%T), roleId=%v (%T), branchId=%v (%T)",
			idValue, idValue, roleIdValue, roleIdValue, branchIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		// ✅ Safely convert branchIdValue to string
		var branchIdStr string
		switch v := branchIdValue.(type) {
		case string:
			branchIdStr = v
		case float64:
			branchIdStr = strconv.Itoa(int(v))
		case int:
			branchIdStr = strconv.Itoa(v)
		default:
			log.Error("❌ Unsupported type for branchId")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid branchId type",
			})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		log.Info("📦 Fetching branch with floors from DB")
		branch, err := settingsService.GetBranchWithFloorsService(dbConnt, branchIdStr)
		if err != nil {
			log.Error("❌ Failed to fetch branch: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch branch details",
			})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Sending response with branch floors\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   branch,
			"token":  token,
		})
	}
}

func GetNewBranchWithFloorWithIdController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("\n\n📥 GetNewBranchWithFloorController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v (%T), roleId=%v (%T), branchId=%v (%T)",
			idValue, idValue, roleIdValue, roleIdValue, branchIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		// ✅ Safely convert branchIdValue to string
		var branchIdStr string
		switch v := branchIdValue.(type) {
		case string:
			branchIdStr = v
		case float64:
			branchIdStr = strconv.Itoa(int(v))
		case int:
			branchIdStr = strconv.Itoa(v)
		default:
			log.Error("❌ Unsupported type for branchId")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid branchId type",
			})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		log.Info("📦 Fetching branch with floors from DB")
		branch, err := settingsService.GetBranchWithFloorsService(dbConnt, branchIdStr)
		if err != nil {
			log.Error("❌ Failed to fetch branch: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch branch details",
			})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Sending response with branch floors\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   branch,
			"token":  token,
		})
	}
}

func UpdateBranchWithFloorController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("\n\nUpdate Branch Controller invoked")

		// Extract JWT context values
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

		// Extract branchId from path param
		paramId := c.Param("id")

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

		userId := 0
		switch v := idValue.(type) {
		case float64:
			userId = int(v)
		case int:
			userId = v
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Invalid user ID type"})
			return
		}

		// Extract role
		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		fmt.Println("Role ID:", roleId)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		// Service call
		err = settingsService.UpdateBranchWithFloor(dbConnt, paramId, &payload.BranchWithFloor, payload.Floors, userId)
		if err != nil {
			log.Error("Failed to update branch with floors: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		log.Info("Branch with Floors and Sections updated successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch updated with Floors and Sections", "token": token})
	}
}

func SoftDeleteBranchController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("\n\nSoft Delete Branch Controller invoked")

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

		// Extract branchId from path param
		paramId := c.Param("id")

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		userId := 0
		switch v := idValue.(type) {
		case float64:
			userId = int(v)
		case int:
			userId = v
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Invalid user ID type"})
			return
		}

		err := settingsService.SoftDeleteBranch(dbConnt, paramId, userId)
		if err != nil {
			log.Error("Failed to soft delete branch: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		log.Info("Branch soft deleted successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Branch soft deleted successfully", "token": token})
	}
}

// ATTRIBUTES
func GetAttributeDataType() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("\n\n📥 GetAttributeDataType invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		log.Info("📦 Fetching all attributes type from DB")
		attributes := settingsService.GetAllAttributesService(dbConnt)
		log.Infof("📊 Attributes fetched: count = %d", len(attributes))

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   attributes,
			"token":  token,
		})
	}
}

func CreateAttributeGroupController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n\n🚀 Create Attribute Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var attributes model.AttributesTable
		if err := c.ShouldBindJSON(&attributes); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Request Body: %+v", attributes)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.CreateAttributesService(dbConnt, &attributes, roleName)
		if err != nil {
			log.Error("❌ Service Error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create category"})
			}
			return
		}

		log.Info("✅ Attribute created successfully\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Category created successfully",
			"token":   token,
		})
	}
}

func GetAttributeGroupController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 GetAttributeGroupController invoked")

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

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		categories := settingsService.GetAttributesService(dbConnt)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   categories,
			"token":  token,
		})
	}
}

func UpdateAttributeGroupController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 UpdateCategoryController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var category model.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			log.Error("❌ Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Request Body: %+v", category)

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("❌ Failed to get role name: " + err.Error())
		} else {
			log.Infof("👤 Role Name: %s", roleName)
		}

		log.Info("🛠️ Calling UpdateCategoryService")
		errH := settingsService.UpdateCategoryService(dbConnt, &category, roleName)
		if errH != nil {
			log.Error("❌ Service error: " + errH.Error())
			if errH.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update category"})
			}
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Category updated successfully\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Category updated successfully",
			"token":   token,
		})
	}
}

func DeleteAttributeGroupController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 DeleteAttributeGroupController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context values (id/roleId/branchId)")
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
			log.Error("❌ Invalid request body or empty category IDs")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid category IDs",
			})
			return
		}
		log.Infof("📦 Bulk delete request: categoryIds=%v, forceDelete=%v", request.CategoryIDs, request.ForceDelete)

		dbConnt, sqlDB := db.InitDB()
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Error("❌ Failed to close DB connection: " + err.Error())
			} else {
				log.Info("✅ DB connection closed")
			}
		}()

		// Step 1: Check for subcategories
		log.Info("🔎 Checking for subcategories in selected categories")
		subcategoriesMap, err := settingsService.CheckSubcategoriesExistence(dbConnt, request.CategoryIDs)
		if err != nil {
			log.Error("❌ Error checking subcategories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Internal server error",
			})
			return
		}

		if len(subcategoriesMap) > 0 && !request.ForceDelete {
			log.Warn("⚠️ Some categories have subcategories. Confirmation needed before force delete.")
			c.JSON(http.StatusConflict, gin.H{
				"status":             false,
				"message":            "Some categories contain subcategories. Deleting them will make subcategories idle.",
				"subcategoriesMap":   subcategoriesMap,
				"confirmationNeeded": true,
			})
			return
		}

		// Step 2: Perform soft delete
		log.Info("🛠️ Proceeding to soft delete categories")

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = settingsService.BulkDeleteCategoriesService(dbConnt, request.CategoryIDs, roleName)
		if err != nil {
			log.Error("❌ Service error during bulk delete: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to delete categories",
			})
			return
		}

		// Optional: you can loop and log each deletion
		log.Infof("✅ Categories soft deleted successfully: %v\n\n", request.CategoryIDs)
		log.Info("\n=================================================================\n")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Categories deleted successfully",
			"token":   token,
		})
	}
}

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

func GetSettingsOverview() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data, err := settingsService.FetchSettingsOverview(dbConn)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   data,
		})
	}
}
