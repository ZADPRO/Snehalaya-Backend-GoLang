package poController

import (
	"net/http"

	poModuleModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/model"
	poService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	roleType "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/GetRoleType"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"
)

func CreatePurchaseOrderController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("🚀 CreatePurchaseOrderController invoked")

		// Token & context validation
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found"})
			return
		}

		var poPayload poModuleModel.PurchaseOrderPayload
		if err := c.ShouldBindJSON(&poPayload); err != nil {
			log.Error("❌ Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		log.Infof("📦 Payload: %+v", poPayload)

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
			roleName = "Unknown"
			log.Warn("⚠️ Role name could not be resolved")
		}

		purchaseOrderNumber, err := poService.CreatePurchaseOrderService(dbConnt, &poPayload, roleName)
		if err != nil {
			log.Error("❌ Service Error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":              true,
			"message":             "Purchase Order created successfully",
			"purchaseOrderNumber": purchaseOrderNumber, // 🧾 send to frontend
			"token":               token,
		})

	}
}

func GetAllPurchaseOrdersController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📥 GetAllPurchaseOrdersController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found"})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data := poService.GetAllPurchaseOrdersService(dbConnt)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "data": data, "token": token})
	}
}

func UpdatePurchaseOrderController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📥 UpdatePurchaseOrderController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found"})
			return
		}

		var poPayload poModuleModel.PurchaseOrderPayload
		if err := c.ShouldBindJSON(&poPayload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, _ := roleType.ExtractIntFromInterface(roleIdValue)
		roleName, _ := roleType.GetRoleTypeNameByID(dbConnt, roleId)

		if err := poService.UpdatePurchaseOrderService(dbConnt, &poPayload, roleName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Purchase Order updated", "token": token})
	}
}
