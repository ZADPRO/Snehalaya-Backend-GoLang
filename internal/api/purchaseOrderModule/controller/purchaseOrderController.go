package purchaseOrderController

import (
	"net/http"

	purchaseOrderModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/model"
	purchaseOrderService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"

)

// CREATE PURCHASE ORDER
func CreatePurchaseOrderController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Create Purchase Order Controller")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")
		createdBy := "Admin"

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		var payload purchaseOrderModel.CreatePORequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Error("Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		err := purchaseOrderService.CreatePurchaseOrderService(dbConnt, &payload, createdBy)
		if err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create purchase order"})
			return
		}

		log.Info("Purchase Order created successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Purchase Order created successfully", "token": token})
	}
}

func GetAllPurchaseOrdersController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Get All Purchase Orders Controller")

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

		purchaseOrders, err := purchaseOrderService.GetAllPurchaseOrdersService(dbConnt)
		if err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch purchase orders",
			})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("Fetched all purchase orders successfully")
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Fetched purchase orders successfully",
			"token":   token,
			"data":    purchaseOrders,
		})
	}
}

func GetPurchaseOrderByIdController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Get All Purchase Orders Controller")

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

		// dbConnt, sqlDB := db.InitDB()
		// defer sqlDB.Close()
		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("Fetched all purchase orders successfully")
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Fetched purchase orders successfully",
			"token":   token,
			// "data":    purchaseOrders,
		})
	}
}

func GetDummyProductsByPOID() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Get Dummy Products by PurchaseOrderId Controller")

		purchaseOrderIdStr := c.Param("purchaseOrderId")
		if purchaseOrderIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Missing purchaseOrderId"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		dummyProducts, err := purchaseOrderService.GetDummyProductsByPOIDService(dbConn, purchaseOrderIdStr)
		if err != nil {
			log.Error("Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to fetch dummy products"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "data": dummyProducts})
	}
}

func UpdateDummyProductStatus() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Update Dummy Product Status Controller")

		var payload struct {
			DummyProductId int         `json:"dummyProductId"`
			Status         interface{} `json:"status"` // can be bool or string
			Reason         string      `json:"reason"` // optional
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Error("Invalid request payload: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := purchaseOrderService.UpdateDummyProductStatusService(dbConn, payload.DummyProductId, payload.Status, payload.Reason)
		if err != nil {
			log.Error("Failed to update dummy product: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Dummy product updated successfully"})
	}
}

// BULK UPDATE - ACCEPT, REJECT, UNDO
func BulkAcceptDummyProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			DummyProductIds []int `json:"dummyProductIds"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := purchaseOrderService.BulkUpdateDummyProducts(dbConn, payload.DummyProductIds, "accept", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Products accepted successfully"})
	}
}

func BulkRejectDummyProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			DummyProductIds []int  `json:"dummyProductIds"`
			Reason          string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := purchaseOrderService.BulkUpdateDummyProducts(dbConn, payload.DummyProductIds, "reject", payload.Reason)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Products rejected successfully"})
	}
}

func BulkUndoDummyProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			DummyProductIds []int `json:"dummyProductIds"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := purchaseOrderService.BulkUpdateDummyProducts(dbConn, payload.DummyProductIds, "undo", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Products reset to pending"})
	}
}
