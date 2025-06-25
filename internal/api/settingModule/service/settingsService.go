package settingsService

import (
	"errors"
	"fmt"
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"

)

// CATEGORIES SERVICE

func CreateCategoryService(db *gorm.DB, category *model.Category) error {
	log := logger.InitLogger()

	// Check for existing category with same name or code and isDelete = false
	var existing model.Category
	err := db.Table("Categories").
		Where(`("categoryName" = ? OR "categoryCode" = ?) AND "isDelete" = ?`, category.CategoryName, category.CategoryCode, false).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate category found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("DB error while checking for duplicates: " + err.Error())
		return err
	}

	// Proceed with creation
	category.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	category.CreatedBy = "Admin"
	return db.Table("Categories").Create(category).Error
}

func GetAllCategoriesService(db *gorm.DB) []model.Category {
	log := logger.InitLogger()
	var categories []model.Category

	err := db.Table("Categories").
		Where(`"isDelete" = ?`, false).
		Find(&categories).Error

	if err != nil {
		log.Error("Get all categories failed: " + err.Error())
	}
	return categories
}

func UpdateCategoryService(db *gorm.DB, category *model.Category) error {
	log := logger.InitLogger()

	// Check for duplicate (excluding current ID)
	var existing model.Category
	err := db.Table("Categories").
		Where(`("categoryName" = ? OR "categoryCode" = ?) AND "refCategoryid" != ? AND "isDelete" = ?`,
			category.CategoryName, category.CategoryCode, category.RefCategoryId, false).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate category found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("DB error while checking for duplicates: " + err.Error())
		return err
	}

	// Proceed with update
	category.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	category.UpdatedBy = "Admin"
	return db.Table("Categories").
		Where(`"refCategoryid" = ?`, category.RefCategoryId).
		Updates(map[string]interface{}{
			"categoryName": category.CategoryName,
			"categoryCode": category.CategoryCode,
			"isActive":     category.IsActive,
		}).Error
}

func GetSubcategoriesByCategory(db *gorm.DB, categoryId string) ([]model.SubCategory, error) {
	var subcategories []model.SubCategory
	err := db.Table("SubCategories").
		Where(`"refCategoryId" = ? AND "isDelete" = false`, categoryId).
		Find(&subcategories).Error

	return subcategories, err
}

func DeleteCategoryService(db *gorm.DB, id string) error {
	log := logger.InitLogger()
	log.Info("Soft deleting category with ID: ", id)

	return db.Table("Categories").
		Where(`"refCategoryid" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": "Admin",
		}).Error
}

// SUB CATEGORIES SERVICE
func CreateSubCategoryService(db *gorm.DB, sub *model.SubCategory) error {
	log := logger.InitLogger()
	log.Info("Inserting SubCategory: ", sub)

	var existing model.SubCategory
	err := db.Table("SubCategories").
		Where(`("subCategoryName" = ? OR "subCategoryCode" = ?) AND "isDelete" = false`, sub.SubCategoryName, sub.SubCategoryCode).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate SubCategory found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("Error checking for duplicate: " + err.Error())
		return err
	}

	sub.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	sub.CreatedBy = "Admin"
	return db.Table("SubCategories").Create(sub).Error
}

func GetAllSubCategoriesService(db *gorm.DB) []model.SubCategory {
	log := logger.InitLogger()
	var subs []model.SubCategory

	err := db.Table("SubCategories").Where(`"isDelete" = false`).Find(&subs).Error
	if err != nil {
		log.Error("Failed to fetch subcategories: " + err.Error())
	}
	return subs
}

func UpdateSubCategoryService(db *gorm.DB, sub *model.SubCategory) error {
	log := logger.InitLogger()
	log.Info("Updating SubCategory ID: ", sub.RefSubCategoryId)

	// 1. Check for duplicates
	var existing model.SubCategory
	err := db.Table("SubCategories").
		Where(`("subCategoryName" = ? OR "subCategoryCode" = ?) AND "refSubCategoryId" != ? AND "isDelete" = false`,
			sub.SubCategoryName, sub.SubCategoryCode, sub.RefSubCategoryId).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate SubCategory found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("Error checking for duplicate: " + err.Error())
		return err
	}

	// 2. Update
	updateData := map[string]interface{}{
		"subCategoryName": sub.SubCategoryName,
		"subCategoryCode": sub.SubCategoryCode,
		"refCategoryId":   sub.RefCategoryId,
		"isActive":        sub.IsActive,
		"updatedAt":       time.Now().Format("2006-01-02 15:04:05"),
		"updatedBY":       "Admin",
	}

	err = db.Table("SubCategories").
		Where(`"refSubCategoryId" = ?`, sub.RefSubCategoryId).
		Updates(updateData).Error

	if err != nil {
		log.Error("Failed to update subcategory: " + err.Error())
		return err
	}

	return nil
}

func DeleteSubCategoryService(db *gorm.DB, id string) error {
	log := logger.InitLogger()
	log.Info("Soft deleting SubCategory with ID: ", id)

	return db.Table("SubCategories").
		Where(`"refSubCategoryId" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBY": "Admin",
		}).Error
}

// BRANCHES SERVICE
func CreateBranchService(db *gorm.DB, branch *model.Branch) error {
	log := logger.InitLogger()
	log.Info("Inserting Branch: ", branch)

	var existing model.Branch

	// Ensure exact column names match your PostgreSQL schema
	err := db.Table(`"Branches"`).
		Where(`("refBranchName" = ? OR "refBranchCode" = ?) AND "isDelete" = false`, branch.RefBranchName, branch.RefBranchCode).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate Branch found")
		return fmt.Errorf("duplicate value found")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Error checking for duplicate: " + err.Error())
		return err
	}

	branch.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	branch.CreatedBy = "Admin"

	// Always use quoted table name to preserve case
	return db.Table(`"Branches"`).Create(branch).Error
}

func GetAllBranchesService(db *gorm.DB) []model.Branch {
	log := logger.InitLogger()
	var branches []model.Branch
	err := db.Table("Branches").Where("isDelete = false").Find(&branches).Error
	if err != nil {
		log.Error("Failed to fetch branches: " + err.Error())
	}
	return branches
}

func UpdateBranchService(db *gorm.DB, branch *model.Branch) error {
	log := logger.InitLogger()
	log.Info("Updating Branch ID: ", branch.RefBranchId)

	var existing model.Branch
	err := db.Table("Branches").
		Where(`("refBranchName" = ? OR "refBranchCode" = ?) AND "refBranchId" != ? AND "isDelete" = false`,
			branch.RefBranchName, branch.RefBranchCode, branch.RefBranchId).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate Branch found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("Error checking for duplicate: " + err.Error())
		return err
	}

	updateData := map[string]interface{}{
		"refBranchName": branch.RefBranchName,
		"refBranchCode": branch.RefBranchCode,
		"refLocation":   branch.RefLocation,
		"refMobile":     branch.RefMobile,
		"refEmail":      branch.RefEmail,
		"refBTId":       branch.RefBTId,
		"isMainBranch":  branch.IsMainBranch,
		"isActive":      branch.IsActive,
		"updatedAt":     time.Now().Format("2006-01-02 15:04:05"),
		"updatedBY":     "Admin",
	}

	return db.Table("Branches").Where("refBranchId = ?", branch.RefBranchId).Updates(updateData).Error
}

func DeleteBranchService(db *gorm.DB, id string) error {
	log := logger.InitLogger()
	log.Info("Soft deleting Branch with ID: ", id)

	return db.Table("Branches").
		Where("refBranchId = ?", id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBY": "Admin",
		}).Error
}
