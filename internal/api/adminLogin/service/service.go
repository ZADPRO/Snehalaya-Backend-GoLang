package service

import (
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/model"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/query"
	transactionLogger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/service"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	becrypt "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Bcrypt"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func AdminLoginService(db *gorm.DB, reqVal model.AdminLoginReq) model.LoginResponse {
	log := logger.InitLogger()

	var AdminLoginModel []model.AdminLoginModelReq

	log.Info("\n\nUser Details -----> \n\n" + reqVal.Username)

	err := db.Raw(query.AdminLoginSQL, reqVal.Username).Scan(&AdminLoginModel).Error
	if err != nil {
		log.Error("Login Service DB Error: " + err.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	if len(AdminLoginModel) == 0 {
		log.Warn("LoginService Invalid Credentials (u) for Username : " + reqVal.Username)
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid Username or Password",
		}
	}

	user := AdminLoginModel[0]
	log.Info("Database query values ---> ", user)

	match := becrypt.ComparePasswords(user.UCDHashedPassword, reqVal.Password)
	log.Warn("\n\nPassword checking => ", match)

	if !match {
		log.Warn("Login Service Invalid Credentials for Username : ", reqVal.Username)
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid Username or Password",
		}
	}

	log.Info("Login service - Logged Successfully for Username : " + reqVal.Username)

	// 🔽 Log the login transaction
	_ = transactionLogger.LogTransaction(
		db,
		user.UserId,
		"Admin", // or user.Username if preferred
		1,       // 1 = Login
		"User logged in: "+reqVal.Username,
	)

	token := accesstoken.CreateToken(user.UserId, user.RoleTypeId, user.UserBranchId)
	log.Info("\n\n\nToken Testing --------->" + token)

	return model.LoginResponse{
		Status:  true,
		Message: "Logged in Successfully",
		User:    &user,
		Token:   token,
	}
}
