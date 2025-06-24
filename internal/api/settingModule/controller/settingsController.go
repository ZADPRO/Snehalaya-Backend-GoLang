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

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		categories := settingsService.GetAllCategoriesService(dbConnt)
		c.JSON(http.StatusOK, gin.H{"status": true, "data": categories})
	}
}

func UpdateCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Update Category Controller")

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

		log.Info("Category updated successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category updated successfully"})
	}
}

func DeleteCategoryController() gin.HandlerFunc {
	log := logger.InitLogger()
	return func(c *gin.Context) {
		log.Info("Delete Category Controller")

		id := c.Param("id")

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := settingsService.DeleteCategoryService(dbConnt, id)
		if err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete category"})
			return
		}

		log.Info("Category deleted successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Category deleted successfully"})
	}
}
