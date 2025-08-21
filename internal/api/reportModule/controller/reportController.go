package reportController

import (
	"net/http"

	reportModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/model"
	reportService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	roleType "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/GetRoleType"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"
)

func GetAllProductsReportController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n\nGetAllProductsReportController invoked")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("\nContext Data: id=%v, roleId=%v, branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		// PAYLOAD VALIDATION
		var productsReportPayload reportModel.ProductsReportPayload
		if err := c.ShouldBindJSON(&productsReportPayload); err != nil {
			log.Error("Invalid pagination payload:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		log.Infof("Request Body: %+v", productsReportPayload)

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		roleId, err := roleType.ExtractIntFromInterface(roleIdValue)
		if err != nil {
			log.Error("Invalid role id:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid role ID",
			})
			return
		}

		roleName, err := roleType.GetRoleTypeNameByID(dbConnt, roleId)
		if err != nil {
			log.Error("Failed to get role type name:", err.Error())
		} else {
			log.Infof("Role name resolved: %s", roleName)
		}

		result, err := reportService.GetAllProductReportsService(dbConnt, &productsReportPayload, roleName)

		if err != nil {
			log.Error("Server Error : " + err.Error())
			return
		}

		log.Info("Reports Fetched Successfully")

		log.Info("\n=================================================================\n")

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Reports fetched successfully",
			"data":    result,
			"token":   token,
		})

	}
}
