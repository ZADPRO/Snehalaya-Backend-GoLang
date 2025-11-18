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
		log.Info("\n\n🚀 CreatePurchaseOrderController invoked")

		// Extract context
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v | roleId=%v | branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found",
			})
			return
		}

		// Bind payload
		var poPayload poModuleModel.PurchaseOrderPayload
		if err := c.ShouldBindJSON(&poPayload); err != nil {
			log.Error("❌ Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		log.Infof("📦 PO Payload Received: %+v", poPayload)

		// DB connection
		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// Role resolve
		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID provided: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Warn("⚠️ Failed to resolve role name: " + err.Error())
			roleName = "Unknown"
		} else {
			log.Infof("👤 Role Name: %s", roleName)
		}

		// Service call
		purchaseOrderNumber, err := poService.CreatePurchaseOrderService(dbConnt, &poPayload, roleName)
		if err != nil {
			log.Error("❌ PO Service Error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		// Create token
		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Infof("✅ Purchase Order created successfully | PO Number: %s", purchaseOrderNumber)
		log.Info("=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":              true,
			"message":             "Purchase Order created successfully",
			"purchaseOrderNumber": purchaseOrderNumber,
			"token":               token,
		})
	}
}

func GetAllPurchaseOrdersController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n📥 GetAllPurchaseOrdersController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v | roleId=%v | branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Missing token claims"})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		log.Info("📡 Fetching all purchase orders...")
		data := poService.GetAllPurchaseOrdersService(dbConnt)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Infof("✅ %d Purchase Orders retrieved\n", len(data))

		c.JSON(http.StatusOK, gin.H{"status": true, "data": data, "token": token})
	}
}

func UpdatePurchaseOrderController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n✏️ UpdatePurchaseOrderController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context: id=%v | roleId=%v | branchId=%v",
			idValue, roleIdValue, branchIdValue,
		)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing claims in token")
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Missing claims"})
			return
		}

		var poPayload poModuleModel.PurchaseOrderPayload
		if err := c.ShouldBindJSON(&poPayload); err != nil {
			log.Error("❌ Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		log.Infof("📦 Update Payload: %+v", poPayload)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, _ := roleType.ExtractIntFromInterface(roleIdValue)
		roleName, _ := roleType.GetRoleTypeNameByID(dbConnt, roleId)

		if err := poService.UpdatePurchaseOrderService(dbConnt, &poPayload, roleName); err != nil {
			log.Error("❌ Update Service Error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Purchase Order updated successfully\n")

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Purchase Order updated", "token": token})
	}
}

func GetAllPurchaseOrdersListController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n📋 GetAllPurchaseOrdersListController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Claims -> id=%v | roleId=%v | branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Missing claims"})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		log.Infof("🗄️ DB Stats: %+v", sqlDB.Stats())

		poList, err := poService.GetAllPurchaseOrdersListService(dbConnt)
		if err != nil {
			log.Error("❌ Failed to fetch Purchase Orders: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		log.Infof("✅ %d Purchase Orders fetched", len(poList))

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{"status": true, "data": poList, "token": token})
	}
}

func UpdatePurchaseOrderProductsController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📝 UpdatePurchaseOrderProductsController invoked")

		var payload []poService.UpdatePOProductRequest // ✅ use struct from service
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Errorf("❌ Invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request payload"})
			return
		}

		if len(payload) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Empty payload"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := poService.UpdatePurchaseOrderProductsService(dbConn, payload)
		if err != nil {
			log.Errorf("❌ Failed to update products: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Database update failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Purchase order products updated successfully",
		})
	}
}

func SavePurchaseOrderProductsController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("🧾 SavePurchaseOrderProductsController invoked")

		var payload poService.SavePurchaseOrderProductsRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Errorf("❌ Invalid request payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request payload"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := poService.SavePurchaseOrderProductsService(dbConn, payload); err != nil {
			log.Errorf("❌ Failed to save PO products: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Database save failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Purchase order products saved successfully",
		})
	}
}

func GetPurchaseOrderDetailsController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📥 GetPurchaseOrderDetailsController invoked")

		// Verify token claims
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found"})
			return
		}

		purchaseOrderNumber := c.Param("purchaseOrderNumber")
		if purchaseOrderNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "purchaseOrderNumber is required"})
			return
		}

		// DB connection
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// Call service
		data, err := poService.GetPurchaseOrderDetailsService(dbConn, purchaseOrderNumber)
		if err != nil {
			log.Error("❌ Failed to fetch PO details: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to fetch purchase order details"})
			return
		}

		// Create new JWT token
		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   data,
			"token":  token,
		})
	}
}

func GetAcceptedProductsController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📦 GetAcceptedProductsController invoked")

		purchaseOrderId := c.Param("purchaseOrderId")
		if purchaseOrderId == "" {
			log.Error("❌ Missing purchaseOrderId in request")
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "purchaseOrderId is required"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		results, err := poService.GetAcceptedProductsService(dbConn, purchaseOrderId)
		if err != nil {
			log.Errorf("❌ Failed to fetch accepted products: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch accepted products",
			})
			return
		}

		log.Infof("✅ Retrieved %d accepted products for PurchaseOrderId: %s", len(results), purchaseOrderId)
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Accepted products retrieved successfully",
			"data":    results,
		})
	}
}

func GetPurchaseOrderDetailsHandler(c *gin.Context) {
	log := logger.InitLogger()
	poNumber := c.Param("purchaseOrderNumber")
	if poNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "PurchaseOrderNumber is required"})
		return
	}

	dbConn, sqlDB := db.InitDB()
	defer sqlDB.Close()

	response, err := poService.GetPurchaseOrderFullDetailsService(dbConn, poNumber)
	if err != nil {
		log.Error("❌ Failed to get PO details: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Purchase order details retrieved successfully",
		"data":    response,
	})
}

func GetAllPurchaseOrderAcceptedProductsController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📥 GetAllPurchaseOrderAcceptedProductsController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "User ID, RoleID, Branch ID not found"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data := poService.GetAllPurchaseOrderAcceptedProductsService(dbConn)

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)
		c.JSON(http.StatusOK, gin.H{"status": true, "data": data, "token": token})
	}
}
