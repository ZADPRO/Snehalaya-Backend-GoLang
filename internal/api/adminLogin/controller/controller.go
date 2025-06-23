package controller

import (
	"net/http"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/model"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"

)

func AdminLoginController() gin.HandlerFunc {

	log := logger.InitLogger()
	return func(c *gin.Context) {
		var reqVal model.AdminLoginReq

		log.Info("\n\nAdmin Login Controller -> \n================")
		// ERROR HANDLING - STATUS CODE IN PARAMS
		if err := c.BindJSON(&reqVal); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Something went wrong, Try again ... " + err.Error(),
			})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.AdminLoginService(dbConnt, reqVal)
		log.Info("Response for controller -> ", resVal)

		response := gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		if resVal.Status {
			response["user"] = resVal.User
			response["token"] = resVal.Token
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}
