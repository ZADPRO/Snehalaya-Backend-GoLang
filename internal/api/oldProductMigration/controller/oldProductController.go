package oldProductController

import (
	"fmt"
	"net/http"

	oldProductMigrationModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/oldProductMigration/model"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

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

func MigrateOldProductsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, roleId, branchId := getUserContext(c)
		if id == nil {
			return
		}

		var oldProdMigr oldProductMigrationModel.MigrateOldProductToDbModel

		if err := c.ShouldBindJSON(&oldProdMigr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		fmt.Print("dbConn", dbConn)
		defer sqlDB.Close()

		token := accesstoken.CreateToken(id, roleId, branchId)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "PO product created successfully", "token": token})

	}
}
