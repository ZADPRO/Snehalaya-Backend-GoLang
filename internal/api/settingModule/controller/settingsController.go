package controller

import (
	"net/http"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/model"
	settingsService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"
)

// CATEGORIES CONTROLLER

func CreateCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Create Category Controller")

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

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		err := settingsService.CreateCategoryService(dbConnt, &category)
		if err != nil {
			log.Error("Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate value found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create category"})
			}
			return
		}

		log.Info("Category created successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category created successfully", "token": token})
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

		err := settingsService.UpdateCategoryService(dbConnt, &category)
		if err != nil {
			log.Error("Service error: " + err.Error())
			if err.Error() == "duplicate value found" {
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
