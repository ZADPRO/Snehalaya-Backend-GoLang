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

type PurchaseOrderController struct{}

func NewPurchaseOrderController() *PurchaseOrderController {
	return &PurchaseOrderController{}
}

func (p *PurchaseOrderController) CreatePurchaseOrderProductsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.InitLogger()
		log.Info("🚀 CreatePurchaseOrderController invoked")

		// Extract context values
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		if !idExists || !roleIdExists || !branchIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Missing user/role/branch context",
			})
			return
		}

		var poPayload poModuleModel.PurchaseOrderProductPayload
		if err := c.ShouldBindJSON(&poPayload); err != nil {
			log.Error("📦 Invalid request body: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}
		log.Infof("📦 Payload: %+v", poPayload)

		db, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("❌ Invalid role ID: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid role ID"})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(db, roleId)
		if err != nil {
			log.Error("🔍 Failed to get role name: " + err.Error())
		} else {
			log.Infof("✅ Role Name resolved: %s", roleName)
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		err = poService.CreatePurchaseOrderProductService(db, &poPayload, roleName)
		if err != nil {
			log.Error("❌ Service error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Purchase Order created successfully",
			"token":   token,
		})
	}
}

func (p *PurchaseOrderController) GetAcceptedPurchaseOrdersController() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.InitLogger()
		log.Info("📦 GetAcceptedPurchaseOrdersController invoked")

		db, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data, err := poService.GetAcceptedPurchaseOrdersService(db)
		if err != nil {
			log.Error("❌ Failed to fetch accepted purchase orders: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to fetch accepted purchase orders",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Accepted Purchase Orders fetched successfully",
			"data":    data,
		})
	}
}
