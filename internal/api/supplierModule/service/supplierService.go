package supplierService

import (
	"fmt"
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreateSupplier(db *gorm.DB, supplier *model.Supplier) error {
	log := logger.InitLogger()

	// Check for duplicates among non-deleted records
	var existing model.Supplier
	err := db.Table("Supplier").
		Where(`"supplierName" = ? AND "supplierCompanyName" = ? AND "supplierCode" = ? AND "isDelete" = false`,
			supplier.SupplierName, supplier.SupplierCompanyName, supplier.SupplierCode).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate supplier found")
		return fmt.Errorf("duplicate supplier with same name, company, and code already exists")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("DB error while checking for duplicate: " + err.Error())
		return err
	}

	// Set metadata
	supplier.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	supplier.CreatedBy = "Admin"
	supplier.IsDelete = false // mark as not deleted

	return db.Table("Supplier").Create(supplier).Error
}

func GetAllSuppliers(db *gorm.DB) ([]model.Supplier, error) {
	var suppliers []model.Supplier
	err := db.Table("Supplier").
		Where(`"isDelete" = false`).
		Find(&suppliers).Error
	return suppliers, err
}

func GetSupplierById(db *gorm.DB, id string) (model.Supplier, error) {
	var supplier model.Supplier
	err := db.Table("Supplier").
		Where(`"supplierId" = ? AND "isDelete" = false`, id).
		First(&supplier).Error
	return supplier, err
}
func UpdateSupplier(db *gorm.DB, supplier *model.Supplier) error {
	supplier.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	supplier.UpdatedBy = "Admin"
	return db.Table("Supplier").
		Where(`"supplierId" = ?`, supplier.SupplierID).
		Updates(supplier).Error
}

func DeleteSupplier(db *gorm.DB, id string) error {
	return db.Table("Supplier").
		Where(`"supplierId" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": "Admin",
		}).Error
}
