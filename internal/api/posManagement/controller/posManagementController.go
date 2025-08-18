package  posManagementController

import (
	"net/http"

	posManagementModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/posManagement/model"
	posManagementService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/posManagement/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"

)

func AddCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, roleId, branchId := getUserContext(c)
		if id == nil {
			return
		}

		var customer posManagementModel.AddCustomer
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := posManagementService.AddCustomer(dbConn, &customer); err != nil {
			c.JSON(http.StatusConflict, gin.H{"status": false, "message": err.Error()})
			return
		}

		token := accesstoken.CreateToken(id, roleId, branchId)
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Customer created successfully",
			"token":   token,
		})
	}
}



// Utility
func getUserContext(c *gin.Context) (interface{}, interface{}, interface{}) {
	idValue, idExists := c.Get("id")
	roleIdValue, roleIdExists := c.Get("roleId")
	branchIdValue, branchIdExists := c.Get("branchId")

	if !idExists || !roleIdExists || !branchIdExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "User ID, RoleID, Branch ID not found in request context.",
		})
		return nil, nil, nil
	}
	return idValue, roleIdValue, branchIdValue
}
