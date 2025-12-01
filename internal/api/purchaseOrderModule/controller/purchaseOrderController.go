package purchaseOrderController

import (
	"net/http"
	"strconv"

	purchaseOrderModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/model"
	purchaseOrderService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	roleType "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/GetRoleType"
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

func GetReceivedDummyProductsController() gin.HandlerFunc {
	return func(c *gin.Context) {

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()
		products, err := purchaseOrderService.GetReceivedDummyProductsService(dbConn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch received products",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   products,
		})
	}
}

func GetReceivedDummyProductsBarcodeController() gin.HandlerFunc {
	return func(c *gin.Context) {

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()
		products, err := purchaseOrderService.GetReceivedDummyProductsBarcodeService(dbConn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch received products",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   products,
		})
	}
}

func CreateProductController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("Create Product Controller")

		id, idOk := c.Get("id")
		roleId, roleOk := c.Get("roleId")
		branchId, branchOk := c.Get("branchId")

		if !idOk || !roleOk || !branchOk {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Missing authentication context",
			})
			return
		}

		var product purchaseOrderModel.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			log.Error("Invalid JSON: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		token := accesstoken.CreateToken(id, roleId, branchId)

		err := purchaseOrderService.CreateProductService(dbConnt, &product)
		if err != nil {
			if err.Error() == "duplicate value found" {
				c.JSON(http.StatusConflict, gin.H{"status": false, "message": "Duplicate SKU found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create product"})
			}
			return
		}

		log.Info("Product created successfully")
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Product created", "token": token})
	}
}
func NewCreatePurchaseOrderController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n\n🚀 Create Purchase Order Controller invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context Data: id=%v, roleId=%v, branchId=%v",
			idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized,
				gin.H{"status": false, "message": "Missing user context"})
			return
		}

		var payload purchaseOrderService.PurchaseOrderPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest,
				gin.H{"status": false, "message": err.Error()})
			return
		}

		log.Infof("📦 PO Payload: %+v", payload)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, _ := roleType.GetRoleTypeNameByID(dbConn, roleId)
		log.Infof("👤 Role Name: %s", roleName)

		createdByFloat, _ := idValue.(float64)
		createdBy := int(createdByFloat)

		result, err := purchaseOrderService.NewCreatePurchaseOrderService(
			dbConn, payload, roleName, createdBy,
		)

		if err != nil {
			log.Error("❌ Service Error: " + err.Error())
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": false, "message": "Failed to create PO"})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Purchase Order created successfully\n\n")
		log.Info("\n=================================================================\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Purchase Order created successfully",
			"data":    result,
			"token":   token,
		})
	}
}

func NewGetAllPurchaseOrdersController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 GetAllPurchaseOrdersController invoked")

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		log.Info("📦 Fetching all purchase orders")
		poList := purchaseOrderService.NewGetAllPurchaseOrdersService(dbConn)

		log.Infof("📊 Purchase Orders fetched: %d", len(poList))

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   poList,
		})
	}
}

func NewGetSinglePurchaseOrderController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📥 GetSinglePurchaseOrderController invoked")

		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		log.Infof("🔍 Fetching PO ID: %d", id)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		result, err := purchaseOrderService.NewGetSinglePurchaseOrderService(dbConn, id)
		if err != nil {
			log.Error("❌ " + err.Error())
			c.JSON(http.StatusNotFound,
				gin.H{"status": false, "message": "PO not found"})
			return
		}

		log.Info("✅ PO fetched successfully")

		c.JSON(http.StatusOK, gin.H{"status": true, "data": result})
	}
}

func NewCreateGRNController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("📦 Create GRN Controller invoked")

		var payload purchaseOrderService.GRNPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Error("❌ Invalid GRN Payload: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		result, err := purchaseOrderService.NewCreateGRNService(dbConn, payload)
		if err != nil {
			log.Error("❌ " + err.Error())
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": false, "message": "Failed to create GRN"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "GRN created successfully",
			"data":    result,
		})
	}
}

func NewGetAllGRNController() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data := purchaseOrderService.NewGetAllGRNService(dbConn)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   data,
		})
	}
}

func NewGetSingleGRNController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data, err := purchaseOrderService.NewGetSingleGRNService(dbConn, id)
		if err != nil {
			c.JSON(http.StatusNotFound,
				gin.H{"status": false, "message": "GRN not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   data,
		})
	}
}

func NewGetInventoryListController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n📦 GetInventoryListController invoked")

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		log.Info("📥 Fetching Inventory List")
		inventoryList, err := purchaseOrderService.NewGetInventoryListService(dbConn)

		if err != nil {
			log.Error("❌ Failed to fetch inventory: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"error":  err.Error(),
			})
			return
		}

		log.Infof("📊 Inventory items fetched: %d", len(inventoryList))

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   inventoryList,
		})
	}
}
