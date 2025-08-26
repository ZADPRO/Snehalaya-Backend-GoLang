package supplierService

import (
	"fmt"
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreateSupplier(db *gorm.DB, supplier *model.Supplier, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ CreateSupplierService invoked")

	log.Infof("📥 Input Supplier: %+v", supplier)
	log.Infof("👤 Created By (roleName): %s", roleName)

	// Check for duplicates
	var existing model.Supplier
	err := db.Table(`"Supplier"`).
		Where(`"supplierName" = ? AND "supplierCompanyName" = ? AND "supplierCode" = ? AND "isDelete" = false`,
			supplier.SupplierName, supplier.SupplierCompanyName, supplier.SupplierCode).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate supplier found")
		return fmt.Errorf("duplicate supplier with same name, company, and code already exists")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("❌ DB error while checking for duplicate: " + err.Error())
		return err
	}

	log.Info("✅ No duplicates found, proceeding to insert")

	// Set metadata
	supplier.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	supplier.CreatedBy = roleName
	supplier.IsDelete = false

	// Insert supplier
	err = db.Table(`"Supplier"`).Create(supplier).Error
	if err != nil {
		log.Error("❌ Failed to insert supplier: " + err.Error())
		return err
	}

	log.Info("✅ Supplier created in DB")

	// Log transaction
	transErr := service.LogTransaction(db, 1, "Admin", 6, "Supplier Created: "+supplier.SupplierName)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func GetAllSuppliers(db *gorm.DB) ([]model.Supplier, error) {
	log := logger.InitLogger()
	log.Info("📘 GetAllSuppliers service invoked")

	var suppliers []model.Supplier
	err := db.Table(`"Supplier"`).
		Where(`"isDelete" = false`).
		Order(`"supplierId" ASC`).
		Find(&suppliers).Error

	if err != nil {
		log.Error("❌ DB Error fetching suppliers: " + err.Error())
	} else {
		log.Infof("📦 Suppliers fetched: %d", len(suppliers))
	}

	return suppliers, err
}

func GetSupplierById(db *gorm.DB, id string) (model.Supplier, error) {
	log := logger.InitLogger()
	log.Infof("📘 GetSupplierById service invoked for ID: %s", id)

	var supplier model.Supplier
	err := db.Table(`"Supplier"`).
		Where(`"supplierId" = ? AND "isDelete" = false`, id).
		First(&supplier).Error

	if err != nil {
		log.Warnf("❌ No supplier found with ID: %s | Error: %v", id, err)
	} else {
		log.Infof("✅ Supplier record found: %+v", supplier)
	}

	return supplier, err
}

func UpdateSupplier(db *gorm.DB, supplier *model.Supplier) error {
	log := logger.InitLogger()
	log.Infof("🔧 UpdateSupplier service invoked for ID: %v", supplier.SupplierID)

	supplier.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	supplier.UpdatedBy = "Admin"

	err := db.Table(`"Supplier"`).
		Where(`"supplierId" = ?`, supplier.SupplierID).
		Updates(supplier).Error

	if err != nil {
		log.Error("❌ Failed to update supplier: " + err.Error())
		return err
	}

	log.Infof("✅ Supplier updated successfully in DB for ID: %v", supplier.SupplierID)

	// Optional: Add transaction log
	transErr := service.LogTransaction(db, 1, "Admin", 2, fmt.Sprintf("Supplier Updated: %s", supplier.SupplierName))
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	}

	return nil
}

func DeleteSupplier(db *gorm.DB, id string) error {
	log := logger.InitLogger()
	log.Infof("🗑️ Soft deleting supplier with ID: %s", id)

	err := db.Table("Supplier").
		Where(`"supplierId" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": "Admin",
		}).Error

	if err != nil {
		log.Error("❌ Error during supplier deletion: " + err.Error())
		return err
	}

	log.Info("✅ Supplier soft-deleted successfully in DB")

	// Optional: Transaction Log
	transErr := service.LogTransaction(db, 1, "Admin", 2, "Supplier Deleted: "+id)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func BulkDeleteSuppliers(db *gorm.DB, ids []int, isDelete bool) error {
	log := logger.InitLogger()
	action := "deleting"
	if !isDelete {
		action = "restoring"
	}
	log.Infof("🗑️ Bulk %s suppliers: %v", action, ids)

	err := db.Table("Supplier").
		Where(`"supplierId" IN (?)`, ids).
		Updates(map[string]interface{}{
			"isDelete":  isDelete,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": "Admin",
		}).Error

	if err != nil {
		log.Error("❌ Error during bulk update: " + err.Error())
		return err
	}

	log.Infof("✅ Suppliers %s successfully in DB", action)

	// Optional: Log transaction for each supplier
	for _, id := range ids {
		transErr := service.LogTransaction(db, 1, "Admin", 2,
			fmt.Sprintf("Supplier %s: %d", action, id))
		if transErr != nil {
			log.Error("⚠️ Failed to log transaction: " + transErr.Error())
		}
	}

	return nil
}
