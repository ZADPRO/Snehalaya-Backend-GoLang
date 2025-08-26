package supplierController

import (
	"fmt"
	"net/http"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/model"
	supplierService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	roleType "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/GetRoleType"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"
)

// SUPPLIER CONTROLLER
func CreateSupplierController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n🚀 Create Supplier Controller invoked")

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

		var supplier model.Supplier
		if err := c.ShouldBindJSON(&supplier); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		log.Infof("📦 Request Body: %+v", supplier)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConn, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
			roleName = "Unknown"
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		err = supplierService.CreateSupplier(dbConn, &supplier, roleName)
		if err != nil {
			log.Error("❌ Service error: " + err.Error())

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

		log.Info("✅ Supplier created successfully")
		log.Info("\n=================================================================\n")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Supplier created successfully",
			"token":   token,
		})
	}
}

func GetAllSuppliersController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n📦 GetAllSuppliersController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data - id: %v, roleId: %v, branchId: %v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing user context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, or Branch ID not found in request context.",
			})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		suppliers, err := supplierService.GetAllSuppliers(dbConn)
		if err != nil {
			log.Error("❌ Failed to fetch suppliers: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch suppliers",
			})
			return
		}

		log.Infof("✅ %d suppliers fetched successfully", len(suppliers))

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   suppliers,
			"token":  token,
		})
	}
}

func GetSupplierByIdController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📦 GetSupplierByIdController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing user context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, or Branch ID not found in request context.",
			})
			return
		}

		id := c.Param("id")
		log.Infof("📌 Request Param: supplierId = %s", id)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		supplier, err := supplierService.GetSupplierById(dbConn, id)
		if err != nil {
			log.Warnf("❌ Supplier not found with ID: %s | Error: %v", id, err)
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Supplier not found"})
			return
		}

		log.Infof("✅ Supplier fetched successfully for ID: %s", id)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   supplier,
			"token":  token,
		})
	}
}

func UpdateSupplierController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("🛠️ UpdateSupplierController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing user context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var supplier model.Supplier
		if err := c.ShouldBindJSON(&supplier); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Supplier Update Data: %+v", supplier)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := supplierService.UpdateSupplier(dbConn, &supplier)
		if err != nil {
			log.Error("❌ Failed to update supplier: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update supplier"})
			return
		}

		log.Infof("✅ Supplier with ID %v updated successfully", supplier.SupplierID)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Supplier updated successfully",
			"token":   token,
		})
	}
}

func DeleteSupplierController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("🗑️ DeleteSupplierController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing user context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		id := c.Param("id")
		log.Infof("🗂️ Supplier ID to delete: %s", id)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := supplierService.DeleteSupplier(dbConn, id)
		if err != nil {
			log.Error("❌ Failed to soft delete supplier: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to delete supplier",
			})
			return
		}

		log.Infof("✅ Supplier with ID %s deleted successfully", id)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Supplier deleted successfully",
			"token":   token,
		})
	}
}

func BulkDeleteSupplierController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("🗑️ BulkDeleteSupplierController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing user context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var req model.BulkDeleteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("❌ Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request payload",
			})
			return
		}

		if len(req.IDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "No supplier IDs provided",
			})
			return
		}

		log.Infof("🔍 Bulk action on suppliers: %v, isDelete=%v", req.IDs, req.IsDelete)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := supplierService.BulkDeleteSuppliers(dbConn, req.IDs, req.IsDelete)
		if err != nil {
			log.Error("❌ Failed to update suppliers: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to update suppliers",
			})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		action := "deleted"
		if !req.IsDelete {
			action = "restored"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": fmt.Sprintf("Suppliers %s successfully", action),
			"token":   token,
		})
	}
}
