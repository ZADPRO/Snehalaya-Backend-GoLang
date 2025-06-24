package settingsService

import (
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreateCategoryService(db *gorm.DB, category *model.Category) error {
	log := logger.InitLogger()
	log.Info("Inserting category: ", category)
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
	log.Info("Updating category: ", category.RefCategoryId)
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
