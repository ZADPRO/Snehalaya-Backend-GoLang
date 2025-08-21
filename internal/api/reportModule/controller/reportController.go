package reportController

import (
	"net/http"

	reportModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/model"
	reportService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	contextutil "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/ExtractUserContext"
	roleType "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/GetRoleType"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"

)

func GetAllProductsReportController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("GetAllProductsReportController invoked")

		ctxUser, ok := contextutil.ExtractUserContext(c)
		if !ok {
			log.Warn("Missing context data")
			return
		}

		log.Infof("Context Data: id=%v, roleId=%v, branchId=%v", ctxUser.ID, ctxUser.RoleID, ctxUser.BranchID)

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

		roleId, err := roleType.ExtractIntFromInterface(ctxUser.RoleID)
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

		token := accesstoken.CreateToken(ctxUser.ID, ctxUser.RoleID, ctxUser.BranchID)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Reports fetched successfully",
			"data":    result,
			"token":   token,
		})
	}
}
