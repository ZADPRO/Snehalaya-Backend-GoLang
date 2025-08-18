package  posManagementService

import (
	// "fmt"
	"time"
	"errors"

	posManagementModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/posManagement/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func AddCustomer(db *gorm.DB, customer *posManagementModel.AddCustomer) error {
	log := logger.InitLogger()

	// Check if mobile already exists
	var existing posManagementModel.AddCustomer
	err := db.Table("customers").
		Where(`"refMobileNo" = ? AND "isDelete" = false`, customer.RefMobileNo).
		First(&existing).Error

	if err == nil {
		// Customer exists
		return errors.New("customer already exists with this mobile number")
	} else if err != gorm.ErrRecordNotFound {
		// DB error
		return err
	}

	// If not found → insert new
	customer.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	customer.CreatedBy = "Admin"
	customer.IsDelete = false

	if err := db.Table("customers").Create(customer).Error; err != nil {
		log.Error("Failed to create customer: " + err.Error())
		return err
	}

	log.Info("New customer created with mobile: " + customer.RefMobileNo)
	return nil
}