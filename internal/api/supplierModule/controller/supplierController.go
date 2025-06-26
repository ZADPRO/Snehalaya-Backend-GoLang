package supplierController

import (
	"net/http"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/model"
	supplierService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

// SUPPLIER CONTROLLER
func CreateSupplierController() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		var supplier model.Supplier
		if err := c.ShouldBindJSON(&supplier); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := supplierService.CreateSupplier(dbConn, &supplier)
		if err != nil {
			if err.Error() == "duplicate supplier with same name, company, and code already exists" {
				c.JSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "Duplicate value found. A supplier with the same name, company, and code already exists.",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to create supplier",
				})
			}
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Supplier created successfully",
			"token":   token,
		})
	}
}

func GetAllSuppliersController() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		suppliers, err := supplierService.GetAllSuppliers(dbConn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to fetch suppliers"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "data": suppliers, "token": token})
	}
}

func GetSupplierByIdController() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		supplier, err := supplierService.GetSupplierById(dbConn, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Supplier not found"})
			return
		}
		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "data": supplier, "token": token})
	}
}

func UpdateSupplierController() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		var supplier model.Supplier
		if err := c.ShouldBindJSON(&supplier); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := supplierService.UpdateSupplier(dbConn, &supplier); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update supplier"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Supplier updated successfully", "token": token})
	}
}

func DeleteSupplierController() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := supplierService.DeleteSupplier(dbConn, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete supplier"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Supplier deleted successfully", "token": token})
	}
}
