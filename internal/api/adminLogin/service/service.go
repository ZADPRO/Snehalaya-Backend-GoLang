package service

import (
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/model"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/query"
	becrypt "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Bcrypt"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func AdminLoginService(db *gorm.DB, reqVal model.AdminLoginReq) model.LoginResponse {
	log := logger.InitLogger()

	var AdminLoginModel []model.AdminLoginModelReq

	log.Info("\n\nUser Details -----> \n\n"+reqVal.Username, reqVal.Password)

	// EXECUTE QUERY WITH USER NAME
	err := db.Raw(query.AdminLoginSQL, reqVal.Username).Scan(&AdminLoginModel).Error

	if err != nil {
		log.Error("Login Service DB Error: " + err.Error())
		return model.LoginResponse{
			Status:  false,
			Message: "Something went wrong, Try Again",
		}
	}

	// CHECK IF ANY USER FOUND OR NOT
	if len(AdminLoginModel) == 0 {
		log.Warn("LoginService Invalid Credentials (u) for Username : " + reqVal.Username)
		return model.LoginResponse{
			Status:  false,
			Message: "Invalid Username or Password",
		}
	}

	// PASSWORD VERIFICATION
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

	return model.LoginResponse{
		Status:   true,
		Message:  "Logged in Successfully",
		RoleType: user.RoleTypeId,
	}
}
