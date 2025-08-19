package settingsService

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	transactionLogger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/model"
	becrypt "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Bcrypt"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	mailService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/MailService"
	"gorm.io/gorm"

)

// CATEGORIES SERVICE

func CreateCategoryService(db *gorm.DB, category *model.Category, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ CreateCategoryService invoked")

	log.Infof("📥 Input Category: %+v", category)
	log.Infof("👤 Created By (roleName): %s", roleName)

	var existing model.Category
	err := db.Table("Categories").
		Where(`("categoryName" = ? OR "categoryCode" = ?) AND "isDelete" = ?`, category.CategoryName, category.CategoryCode, false).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate category found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("❌ DB error during duplicate check: " + err.Error())
		return err
	}

	log.Info("✅ No duplicates found, proceeding to create category")

	category.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	category.CreatedBy = roleName

	err = db.Table("Categories").Create(category).Error
	if err != nil {
		log.Error("❌ Failed to insert category: " + err.Error())
		return err
	}

	log.Info("✅ Category created in DB, logging transaction...")

	// Transaction Logging
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 2, "Category Created: "+category.CategoryName)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func GetAllCategoriesService(db *gorm.DB) []model.Category {
	log := logger.InitLogger()
	log.Info("🛠️ GetAllCategoriesService invoked")

	var categories []model.Category

	err := db.Table("Categories").
		Where(`"isDelete" = ?`, false).
		Order(`"refCategoryid" ASC`).
		Find(&categories).Error

	if err != nil {
		log.Error("❌ Failed to fetch categories: " + err.Error())
	} else {
		log.Infof("✅ Retrieved %d categories from DB", len(categories))
	}

	return categories
}

func UpdateCategoryService(db *gorm.DB, category *model.Category, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ UpdateCategoryService invoked")

	log.Infof("📥 Category to update: %+v", category)
	log.Infof("👤 Updated By: %s", roleName)

	// Check for duplicate category (excluding the current one)
	var existing model.Category
	err := db.Table("Categories").
		Where(`("categoryName" = ? OR "categoryCode" = ?) AND "refCategoryid" != ? AND "isDelete" = ?`,
			category.CategoryName, category.CategoryCode, category.RefCategoryId, false).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate category found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("❌ DB error while checking for duplicates: " + err.Error())
		return err
	}

	log.Info("✅ No duplicates found, proceeding with update")

	// Update metadata
	category.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	category.UpdatedBy = roleName

	// Log transaction
	log.Info("📝 Logging transaction for category update")
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 3, "Category Updated: "+category.CategoryName)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	// Perform DB update
	log.Info("🔧 Updating category in DB")
	err = db.Table("Categories").
		Where(`"refCategoryid" = ?`, category.RefCategoryId).
		Updates(map[string]interface{}{
			"categoryName": category.CategoryName,
			"categoryCode": category.CategoryCode,
			"isActive":     category.IsActive,
		}).Error

	if err != nil {
		log.Error("❌ Category update failed: " + err.Error())
	} else {
		log.Info("✅ Category updated in DB")
	}

	return err
}

func GetSubcategoriesByCategory(db *gorm.DB, categoryId string) ([]model.SubCategory, error) {
	log := logger.InitLogger()
	log.Infof("🔎 Fetching subcategories for category ID: %s", categoryId)

	var subcategories []model.SubCategory
	err := db.Table("SubCategories").
		Where(`"refCategoryId" = ? AND "isDelete" = false`, categoryId).
		Find(&subcategories).Error

	if err != nil {
		log.Error("❌ Failed to fetch subcategories: " + err.Error())
	} else {
		log.Infof("📊 Found %d subcategories", len(subcategories))
	}

	return subcategories, err
}

func DeleteCategoryService(db *gorm.DB, id string) error {
	log := logger.InitLogger()
	log.Infof("🗑️ Soft deleting category with ID: %s", id)

	err := db.Table("Categories").
		Where(`"refCategoryid" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": "Admin", // Optional: you can pass dynamic roleName if needed
		}).Error

	if err != nil {
		log.Error("❌ Failed to soft delete category: " + err.Error())
		return err
	}

	log.Info("✅ Category soft deleted successfully")

	// Log transaction
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 4, "Category Deleted: "+id)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func BulkDeleteCategoriesService(db *gorm.DB, ids []int, roleName string) error {
	log := logger.InitLogger()
	log.Infof("🗑️ Soft deleting categories with IDs: %v", ids)

	err := db.Table("Categories").
		Where(`"refCategoryid" IN (?)`, ids).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": roleName,
		}).Error

	if err != nil {
		log.Error("❌ Failed to soft delete categories: " + err.Error())
		return err
	}

	log.Info("✅ Bulk category soft delete successful")

	// Optional: Log one transaction for all deletions
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 5, fmt.Sprintf("Bulk Category Delete: %v", ids))
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func CheckSubcategoriesExistence(db *gorm.DB, categoryIDs []int) (map[string][]model.SubCategory, error) {
	log := logger.InitLogger()
	log.Infof("🔎 Checking subcategory existence for category IDs: %v", categoryIDs)

	var subcategories []model.SubCategory
	result := make(map[string][]model.SubCategory)

	err := db.Table("SubCategories").
		Where(`"refCategoryId" IN (?) AND "isDelete" = false`, categoryIDs).
		Find(&subcategories).Error
	if err != nil {
		log.Error("❌ Error fetching subcategories: " + err.Error())
		return nil, err
	}

	for _, sub := range subcategories {
		key := strconv.Itoa(sub.RefCategoryId)
		result[key] = append(result[key], sub)
	}

	log.Infof("📊 Found %d subcategories under %d categories", len(subcategories), len(result))

	return result, nil
}

// SUB CATEGORIES SERVICE
func CreateSubCategoryService(db *gorm.DB, sub *model.SubCategory, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ CreateSubCategoryService invoked")

	log.Infof("📥 Input SubCategory: %+v", sub)
	log.Infof("👤 Created By (roleName): %s", roleName)

	var existing model.SubCategory
	err := db.Table("SubCategories").
		Where(`("subCategoryName" = ? OR "subCategoryCode" = ?) AND "isDelete" = false`,
			sub.SubCategoryName, sub.SubCategoryCode).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate SubCategory found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("❌ DB error during duplicate check: " + err.Error())
		return err
	}

	log.Info("✅ No duplicates found, proceeding to create subcategory")

	sub.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	sub.CreatedBy = roleName

	err = db.Table("SubCategories").Create(sub).Error
	if err != nil {
		log.Error("❌ Failed to insert subcategory: " + err.Error())
		return err
	}

	log.Info("✅ SubCategory created in DB, logging transaction...")

	// ✅ Transaction Logging
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 2, "SubCategory Created: "+sub.SubCategoryName)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func GetAllSubCategoriesService(db *gorm.DB) []model.SubCategoryResponse {
	log := logger.InitLogger()
	var subs []model.SubCategoryResponse

	err := db.Table(`"SubCategories" AS sub`).
		Select(`sub."refSubCategoryId" AS "refSubCategoryId",
	        sub."subCategoryName" AS "subCategoryName",
	        sub."refCategoryId" AS "refCategoryId",
	        cat."categoryName" AS "categoryName",
	        sub."subCategoryCode" AS "subCategoryCode",
	        sub."isActive" AS "isActive",
	        sub."createdAt" AS "createdAt",
	        sub."createdBy" AS "createdBy",
	        sub."updatedAt" AS "updatedAt",
	        sub."updatedBY" AS "updatedBY"`).
		Joins(`JOIN "Categories" AS cat ON sub."refCategoryId" = cat."refCategoryid"`).
		Where(`sub."isDelete" = false AND cat."isDelete" = false`).
		Order(`sub."refSubCategoryId" ASC`).
		Scan(&subs).Error

	fmt.Printf("Result: %+v\n", subs)

	if err != nil {
		log.Error("Failed to fetch subcategories with category names: " + err.Error())
	}
	return subs
}

func UpdateSubCategoryService(db *gorm.DB, sub *model.SubCategory, roleName string) error {
	log := logger.InitLogger()
	log.Infof("🛠️ UpdateSubCategoryService invoked for ID: %d", sub.RefSubCategoryId)

	log.Infof("📥 Input SubCategory: %+v", sub)
	log.Infof("👤 Updated By (roleName): %s", roleName)

	// 1. Check for duplicates
	var existing model.SubCategory
	err := db.Table("SubCategories").
		Where(`("subCategoryName" = ? OR "subCategoryCode" = ?) AND "refSubCategoryId" != ? AND "isDelete" = false`,
			sub.SubCategoryName, sub.SubCategoryCode, sub.RefSubCategoryId).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate SubCategory found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("❌ Error checking for duplicate: " + err.Error())
		return err
	}

	// 2. Perform update
	updateData := map[string]interface{}{
		"subCategoryName": sub.SubCategoryName,
		"subCategoryCode": sub.SubCategoryCode,
		"refCategoryId":   sub.RefCategoryId,
		"isActive":        sub.IsActive,
		"updatedAt":       time.Now().Format("2006-01-02 15:04:05"),
		"updatedBY":       roleName,
	}

	err = db.Table("SubCategories").
		Where(`"refSubCategoryId" = ?`, sub.RefSubCategoryId).
		Updates(updateData).Error

	if err != nil {
		log.Error("❌ Failed to update SubCategory: " + err.Error())
		return err
	}

	log.Info("✅ SubCategory updated in DB successfully")
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

func BulkDeleteSubCategoriesService(db *gorm.DB, ids []int, roleName string) error {
	log := logger.InitLogger()
	log.Infof("🗑️ Soft deleting subcategories with IDs: %v", ids)

	err := db.Table("SubCategories").
		Where(`"refSubCategoryId" IN (?)`, ids).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBY": roleName,
		}).Error

	if err != nil {
		log.Error("❌ Failed to soft delete subcategories: " + err.Error())
		return err
	}

	log.Info("✅ Bulk subcategory soft delete successful")

	// Log one transaction entry (optional)
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 6, fmt.Sprintf("Bulk SubCategory Delete: %v", ids))
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

// BRANCHES SERVICE
func CreateBranchService(db *gorm.DB, branch *model.Branch, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ CreateBranchService invoked")

	log.Infof("📥 Input Branch: %+v", branch)
	log.Infof("👤 Created By (roleName): %s", roleName)

	var existing model.Branch
	err := db.Table(`"Branches"`).
		Where(`("refBranchName" = ? OR "refBranchCode" = ?) AND "isDelete" = false`, branch.RefBranchName, branch.RefBranchCode).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate branch found")
		return fmt.Errorf("duplicate value found")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("❌ DB error during duplicate check: " + err.Error())
		return err
	}

	log.Info("✅ No duplicates found, proceeding to create branch")

	branch.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	branch.CreatedBy = roleName

	err = db.Table(`"Branches"`).Create(branch).Error
	if err != nil {
		log.Error("❌ Failed to insert branch: " + err.Error())
		return err
	}

	log.Info("✅ Branch created in DB, logging transaction...")

	// Transaction Logging (TransTypeID = 4 for branch)
	transErr := transactionLogger.LogTransaction(db, 1, roleName, 4, "Branch Created: "+branch.RefBranchName)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func GetAllBranchesService(db *gorm.DB) ([]model.Branch, error) {
	log := logger.InitLogger()
	var branches []model.Branch

	log.Info("🔎 Fetching all active branches sorted by refBranchId ASC")

	err := db.Table(`"Branches"`).
		Where(`"isDelete" = false`).
		Order(`"refBranchId" ASC`).
		Find(&branches).Error

	if err != nil {
		log.Error("❌ Failed to fetch branches: " + err.Error())
		return nil, err
	}

	log.Infof("✅ %d branches fetched successfully", len(branches))
	return branches, nil
}

func UpdateBranchService(db *gorm.DB, branch *model.Branch, roleName string) error {
	log := logger.InitLogger()
	log.Infof("🔧 UpdateBranchService invoked for Branch ID: %d", branch.RefBranchId)

	var existing model.Branch
	err := db.Table(`"Branches"`).
		Where(`("refBranchName" = ? OR "refBranchCode" = ?) AND "refBranchId" != ? AND "isDelete" = false`,
			branch.RefBranchName, branch.RefBranchCode, branch.RefBranchId).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate branch name or code found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("❌ Error checking for duplicates: " + err.Error())
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
		"updatedBy":     roleName,
	}

	err = db.Table(`"Branches"`).
		Where(`"refBranchId" = ?`, branch.RefBranchId).
		Updates(updateData).Error

	if err != nil {
		log.Error("❌ Failed to update branch: " + err.Error())
		return err
	}

	log.Info("✅ Branch updated successfully in DB")

	// Log transaction history
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 3, fmt.Sprintf("Branch Updated: %s", branch.RefBranchName))
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func DeleteBranchService(db *gorm.DB, id string, roleName string) error {
	log := logger.InitLogger()
	log.Infof("🛠️ DeleteBranchService invoked for Branch ID: %s by %s", id, roleName)

	// Fetch the branch before deleting for logging purpose
	var branch model.Branch
	err := db.Table(`"Branches"`).
		Where(`"refBranchId" = ? AND "isDelete" = false`, id).
		First(&branch).Error
	if err != nil {
		log.Error("❌ Failed to fetch branch: " + err.Error())
		return err
	}

	// Perform soft delete
	err = db.Table(`"Branches"`).
		Where(`"refBranchId" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": roleName,
		}).Error

	if err != nil {
		log.Error("❌ Failed to soft delete branch: " + err.Error())
		return err
	}

	log.Info("✅ Branch soft deleted in DB")

	// Transaction logging
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 4, fmt.Sprintf("Branch Deleted: %s", branch.RefBranchName))
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

// BRANCH SERVICE WITH FLOOR
func CreateNewBranchWithFloor(db *gorm.DB, branch *model.BranchWithFloor, floors []struct {
	FloorName string
	FloorCode string
	Sections  []struct {
		CategoryId       int
		RefSubCategoryId int
		SectionName      string
		SectionCode      string
	}
}, userId int) error {
	log := logger.InitLogger()
	log.Info("Creating new branch: " + branch.RefBranchName)

	// Duplicate check
	var existing model.BranchWithFloor
	err := db.Table(`"Branches"`).Where(`("refBranchName" = ? OR "refBranchCode" = ?) AND "isDelete" = false`, branch.RefBranchName, branch.RefBranchCode).First(&existing).Error
	if err == nil {
		return fmt.Errorf("duplicate branch found")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := db.Begin()
	branch.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	branch.CreatedBy = "Admin"
	if err := tx.Table(`"Branches"`).Create(branch).Error; err != nil {
		tx.Rollback()
		return err
	}
	log.Info("Inserted Branch ID: %d", branch.RefBranchId)

	// Insert floors and sections
	for _, floor := range floors {
		floorModel := model.Floors{
			RefBranchId:  branch.RefBranchId,
			RefFloorName: floor.FloorName,
			RefFloorCode: floor.FloorCode,
			IsActive:     "true",
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:    "Admin",
		}
		if err := tx.Table(`"refFloors"`).Create(&floorModel).Error; err != nil {
			tx.Rollback()
			return err
		}
		log.Info("Inserted Floor ID: %d", floorModel.RefFloorId)

		for _, section := range floor.Sections {
			sectionModel := model.Sections{
				RefFloorId:       floorModel.RefFloorId,
				RefSectionName:   section.SectionName,
				RefSectionCode:   section.SectionCode,
				RefCategoryId:    section.CategoryId,
				RefSubCategoryId: section.RefSubCategoryId,
				IsActive:         "true",
				CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:        "Admin",
			}
			if err := tx.Table(`"refSections"`).Create(&sectionModel).Error; err != nil {
				tx.Rollback()
				return err
			}
			log.Info("Inserted Section ID: %d", sectionModel.RefSectionId)
		}
	}

	// Insert Transaction History
	history := model.TransactionHistory{
		RefTransTypeId:  1,
		RefTransHisData: fmt.Sprintf("Branch Created: %s with Floors and Sections", branch.RefBranchName),
		CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:       "Admin",
		RefUserId:       userId,
	}
	if err := tx.Table(`"TransactionHistory"`).Create(&history).Error; err != nil {
		tx.Rollback()
		return err
	}

	log.Info("Inserted transaction history for Branch ID: %d", branch.RefBranchId)
	tx.Commit()
	return nil
}

func GetBranchWithFloorsService(db *gorm.DB, branchIdStr string) ([]model.BranchResponse, error) {
	log := logger.InitLogger()
	log.Info("🛠️ GetBranchWithFloorsService invoked")

	var rows []struct {
		RefBranchId      int    `gorm:"column:ref_branch_id"`
		RefBranchName    string `gorm:"column:ref_branch_name"`
		RefBranchCode    string `gorm:"column:ref_branch_code"`
		RefLocation      string `gorm:"column:ref_location"`
		RefMobile        string `gorm:"column:ref_mobile"`
		RefEmail         string `gorm:"column:ref_email"`
		IsMainBranch     bool   `gorm:"column:is_main_branch"`
		IsActive         bool   `gorm:"column:is_active"`
		IsOnline         bool   `gorm:"column:is_online"`
		IsOffline        bool   `gorm:"column:is_offline"`
		RefFloorId       int    `gorm:"column:ref_floor_id"`
		RefFloorName     string `gorm:"column:ref_floor_name"`
		RefFloorCode     string `gorm:"column:ref_floor_code"`
		RefSectionId     int    `gorm:"column:ref_section_id"`
		RefSectionName   string `gorm:"column:ref_section_name"`
		RefSectionCode   string `gorm:"column:ref_section_code"`
		RefCategoryId    int    `gorm:"column:ref_category_id"`
		RefSubCategoryId int    `gorm:"column:ref_sub_category_id"`
	}

	query := `
		SELECT 
			br."refBranchId"      AS "ref_branch_id",
			br."refBranchName"    AS "ref_branch_name",
			br."refBranchCode"    AS "ref_branch_code",
			br."refLocation"      AS "ref_location",
			br."refMobile"        AS "ref_mobile",
			br."refEmail"         AS "ref_email",
			br."isMainBranch"     AS "is_main_branch",
			br."isActive"         AS "is_active",
			br."isOnline"         AS "is_online",
			br."isOffline"        AS "is_offline",
			rf."refFloorId"       AS "ref_floor_id",
			rf."refFloorName"     AS "ref_floor_name",
			rf."refFloorCode"     AS "ref_floor_code",
			rs."refSectionId"     AS "ref_section_id",
			rs."refSectionName"   AS "ref_section_name",
			rs."refSectionCode"   AS "ref_section_code",
			rs."refCategoryId"    AS "ref_category_id",
			rs."refSubCategoryId" AS "ref_sub_category_id"
		FROM "Branches" br
		LEFT JOIN "refFloors" rf ON br."refBranchId" = rf."refBranchId"
		LEFT JOIN "refSections" rs ON rf."refFloorId" = rs."refFloorId"
		ORDER BY br."refBranchId" ASC;
	`

	log.Infof("📦 Executing query to fetch branches with floors and sections")
	err := db.Raw(query).Scan(&rows).Error

	if err != nil {
		log.Errorf("❌ DB query failed: %v", err)
		return nil, err
	}
	log.Infof("✅ Query executed successfully, fetched %d rows", len(rows))

	if len(rows) == 0 {
		log.Warn("⚠️ No branches found")
		return nil, errors.New("no branches found")
	}

	// Map to group branches
	branchMap := make(map[int]*model.BranchResponse)
	for _, r := range rows {
		// Create branch if not exists
		if _, exists := branchMap[r.RefBranchId]; !exists {
			log.Infof("🏢 Creating new branch: ID=%d, Name=%s", r.RefBranchId, r.RefBranchName)
			branch := model.BranchResponse{
				RefBranchId:   r.RefBranchId,
				RefBranchName: r.RefBranchName,
				RefBranchCode: r.RefBranchCode,
				RefLocation:   r.RefLocation,
				RefMobile:     r.RefMobile,
				RefEmail:      r.RefEmail,
				IsMainBranch:  r.IsMainBranch,
				IsActive:      r.IsActive,
				IsOnline:      r.IsOnline,
				IsOffline:     r.IsOffline,
				Floors:        []model.FloorResponse{},
			}
			branchMap[r.RefBranchId] = &branch
		}

		branch := branchMap[r.RefBranchId]

		// Add floor if exists
		if r.RefFloorId != 0 {
			var floor *model.FloorResponse
			for i := range branch.Floors {
				if branch.Floors[i].RefFloorId == r.RefFloorId {
					floor = &branch.Floors[i]
					break
				}
			}
			if floor == nil {
				log.Infof("  🏬 Adding new floor: ID=%d, Name=%s (BranchID=%d)", r.RefFloorId, r.RefFloorName, r.RefBranchId)
				branch.Floors = append(branch.Floors, model.FloorResponse{
					RefFloorId: r.RefFloorId,
					FloorName:  r.RefFloorName,
					FloorCode:  r.RefFloorCode,
					Sections:   []model.SectionResponse{},
				})
				floor = &branch.Floors[len(branch.Floors)-1]
			}

			// Add section if exists
			if r.RefSectionId != 0 {
				log.Infof("    📂 Adding section: ID=%d, Name=%s (FloorID=%d)", r.RefSectionId, r.RefSectionName, r.RefFloorId)
				floor.Sections = append(floor.Sections, model.SectionResponse{
					RefSectionId:     r.RefSectionId,
					SectionName:      r.RefSectionName,
					SectionCode:      r.RefSectionCode,
					CategoryId:       r.RefCategoryId,
					RefSubCategoryId: r.RefSubCategoryId,
				})
			} else {
				log.Infof("⚠️ No section found for FloorID=%d", r.RefFloorId)
			}
		} else {
			log.Infof("⚠️ No floor found for BranchID=%d", r.RefBranchId)
		}
	}

	// Convert map to slice
	var response []model.BranchResponse
	for _, b := range branchMap {
		response = append(response, *b)
	}

	log.Infof("✅ Successfully built response with %d branches", len(response))
	return response, nil
}

// ATTRIBUTES
func GetAllAttributesService(db *gorm.DB) []model.AttributeGroupTable {
	log := logger.InitLogger()
	log.Info("🛠️ GetAllAttributesService invoked")

	var attributes []model.AttributeGroupTable

	err := db.Table(`"AttributeGroup"`).
		Order(`"attributeGroupId" ASC`).
		Find(&attributes).Error

	if err != nil {
		log.Error("❌ Failed to fetch attributes: " + err.Error())
		return []model.AttributeGroupTable{}
	}

	log.Infof("✅ Retrieved %d attributes from DB", len(attributes))
	return attributes
}

func CreateAttributesService(db *gorm.DB, attribute *model.AttributesTable, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ CreateAttributesService invoked")

	log.Infof("📥 Input Attribute: %+v", attribute)
	log.Infof("👤 Created By (roleName): %s", roleName)

	var existing model.AttributesTable
	err := db.Table("Attributes").
		Where(`("attributeGroupId" = ? AND "attributeValue" = ?) AND "isDelete" = ?`, attribute.AttributeGroupId, attribute.AttributeValue, false).
		First(&existing).Error

	if err == nil {
		log.Warn("⚠️ Duplicate Attribute found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("❌ DB error during duplicate check: " + err.Error())
		return err
	}

	log.Info("✅ No duplicates found, proceeding to create Attribute")

	attribute.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	attribute.CreatedBy = roleName

	err = db.Table("Attributes").Create(attribute).Error
	if err != nil {
		log.Error("❌ Failed to insert Attribute: " + err.Error())
		return err
	}

	log.Info("✅ Attributes created in DB, logging transaction...")

	// Transaction Logging
	transErr := transactionLogger.LogTransaction(db, 1, "Admin", 2, "Attribute Created: "+attribute.AttributeValue)
	if transErr != nil {
		log.Error("⚠️ Failed to log transaction: " + transErr.Error())
	} else {
		log.Info("📘 Transaction log saved successfully")
	}

	return nil
}

func GetAttributesService(db *gorm.DB) []model.AttributeWithGroup {
	log := logger.InitLogger()
	log.Info("🛠️ GetAttributesService invoked")

	query := `
		SELECT 
		a."attributeId",
		a."attributeGroupId",
		ag."attributeGroupName",
		a."attributeKey",
		a."attributeValue",
		a."createdAt",
		a."createdBy",
		a."updatedAt",
		a."updatedBy",
		a."isDelete"
	FROM "Attributes" AS a
	LEFT JOIN "AttributeGroup" AS ag
		ON a."attributeGroupId" = ag."attributeGroupId";
	`

	var attributes []model.AttributeWithGroup

	err := db.Raw(query).Scan(&attributes).Error
	if err != nil {
		log.Error("❌ Failed to fetch Attributes: " + err.Error())
		return []model.AttributeWithGroup{}
	}

	log.Infof("✅ Retrieved %d Attributes from DB", len(attributes))
	return attributes
}

// EMPLOYEE - SELECT ROLE TYPE

func GetUserRoleTypeService(db *gorm.DB) []model.RoleType {
	log := logger.InitLogger()
	var roleTypes []model.RoleType

	err := db.Find(&roleTypes).Error
	if err != nil {
		log.Error("Error fetching role types: ", err)
		return nil
	}

	return roleTypes
}

func CreateEmployeeService(db *gorm.DB, data *model.EmployeePayload) error {
	txn := db.Begin()
	if txn.Error != nil {
		return txn.Error
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	createdBy := "Admin"

	// 🔍 Step 1: Duplicate check on username, email, or mobile
	var existingCount int64
	if err := txn.Table(`"refUserAuthCred"`).
		Where(`"refUACUsername" = ?`, data.Username).
		Or(`"refUserId" IN (
		SELECT "refUserId" FROM "refUserCommunicationDetails"
		WHERE "refUCDEmail" = ? OR "refUCDMobile" = ?
	)`, data.Email, data.Mobile).
		Count(&existingCount).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("error checking for duplicates: %w", err)
	}

	if existingCount > 0 {
		txn.Rollback()
		return fmt.Errorf("user with same username/email/mobile already exists")
	}

	// 🔤 Step 2: Generate refUserCustId
	roleAbbreviations := map[int]string{
		1:  "SUP", // Super Admin
		2:  "ADM", // Admin
		3:  "ACC", // Accounts Manager
		4:  "STM", // Store Manager
		5:  "PMN", // Purchase Manager
		6:  "BEX", // Billing Executive
		7:  "SAL", // Sales Executive
		8:  "SEO", // SEO
		9:  "CSP", // Customer Support
		10: "SUP", // Supplier
	}

	abbr, ok := roleAbbreviations[data.RoleTypeId]
	if !ok {
		txn.Rollback()
		return fmt.Errorf("invalid role type ID: %d", data.RoleTypeId)
	}

	var userCount int64
	if err := txn.Table(`"Users"`).Where(`"refRTId" = ?`, data.RoleTypeId).Count(&userCount).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("error counting users for role: %w", err)
	}

	refUserCustId := fmt.Sprintf("Z%02d%s%04d", data.RoleTypeId, abbr, userCount+1)

	// 👤 Step 3: Insert into Users table
	user := model.User{
		RefRTId:            data.RoleTypeId,
		RefUserCustId:      refUserCustId,
		RefUserFName:       data.FirstName,
		RefUserLName:       data.LastName,
		RefUserDesignation: data.Designation,
		RefUserStatus:      map[bool]string{true: "Active", false: "In Active"}[data.RefUserStatus],
		RefUserBranchId:    data.RefUserBranchId,
		CreatedAt:          timestamp,
		CreatedBy:          createdBy,
	}
	if err := txn.Table(`"Users"`).Create(&user).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Step 3.5: Generate password and hash
	if len(data.Username) < 4 || len(data.Mobile) < 4 {
		txn.Rollback()
		return fmt.Errorf("username or mobile too short to generate password")
	}

	password := data.Username[:4] + data.Mobile[len(data.Mobile)-4:]
	hashedPassword, err := becrypt.HashPassword(password)
	if err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 🔐 Step 4: Insert into refUserAuthCred
	auth := model.UserAuth{
		RefUserId:            user.RefUserId,
		RefUACUsername:       data.Username,
		RefUACPassword:       password,       // store raw password if needed (not recommended in prod)
		RefUACHashedPassword: hashedPassword, // actual used password
		CreatedAt:            timestamp,
		CreatedBy:            createdBy,
	}

	if err := txn.Table(`"refUserAuthCred"`).Create(&auth).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to create auth: %w", err)
	}

	// 📞 Step 5: Insert into refUserCommunicationDetails
	comm := model.UserCommunication{
		RefUserId:    user.RefUserId,
		RefUCDEmail:  data.Email,
		RefUCDMobile: data.Mobile,
		RefUCDDoorNo: data.DoorNumber,
		RefUCDStreet: data.StreetName,
		RefUCDCity:   data.City,
		RefUCDState:  data.State,
		CreatedAt:    timestamp,
		CreatedBy:    createdBy,
	}
	if err := txn.Table(`"refUserCommunicationDetails"`).Create(&comm).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to create communication: %w", err)
	}

	emailBody := fmt.Sprintf(`
		<table width="100%%" cellspacing="0" cellpadding="0" style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 5px;">
			<tr>
				<td align="center">
					<table width="600" cellspacing="0" cellpadding="0" style="background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
						<tr>
							<td style="background-color: #8B0000; padding: 20px; text-align: center; color: white;">
								<h2 style="margin: 0;">🎉 Welcome to Snehalayaa Silks Family! 🎉</h2>
							</td>
						</tr>
						<tr>
							<td style="padding: 30px; color: #333;">
								<p>Dear <strong>%s %s</strong>,</p>

								<p>We are thrilled to welcome you onboard as a valued member of our ERP project at <strong>Snehalayaa Silks</strong>.</p>

								<p>Your presence marks the beginning of a new chapter in our journey towards excellence in textile innovation and digital transformation.</p>

								<h3 style="color: #8B0000; border-bottom: 1px solid #ccc; padding-bottom: 5px; text-align: center;">Your Employee Credentials</h3>
								<table width="80%%" cellpadding="10" cellspacing="0" style="border-collapse: collapse; margin: 10px auto;">
									<tr style="background-color: #f2f2f2;">
										<td width="40%%" style="border: 1px solid #ddd;"><strong>Employee ID</strong></td>
										<td style="border: 1px solid #ddd;">%s</td>
									</tr>
									<tr>
										<td style="border: 1px solid #ddd;"><strong>Username</strong></td>
										<td style="border: 1px solid #ddd;">%s</td>
									</tr>
									<tr style="background-color: #f2f2f2;">
										<td style="border: 1px solid #ddd;"><strong>Temporary Password</strong></td>
										<td style="border: 1px solid #ddd;">%s</td>
									</tr>
								</table>

								<p style="margin-top: 20px;">Please log in to the system and update your password at the earliest for security purposes.</p>

								<hr style="margin: 30px 0; border: none; border-top: 1px solid #ccc;" />

								<p style="font-style: italic; color: #555;">Together, let's weave success into every thread of Snehalayaa Silks!</p>

								<p style="margin-top: 30px;">Warm regards,</p>
								<p><strong>HR Team</strong><br/>Snehalayaa Silks ERP Project</p>
							</td>
						</tr>
						<tr>
							<td style="background-color: #f2f2f2; text-align: center; padding: 15px; font-size: 12px; color: #999;">
								© 2025 Snehalayaa Silks. All rights reserved.
							</td>
						</tr>
					</table>
				</td>
			</tr>
		</table>
		`, data.FirstName, data.LastName, refUserCustId, data.Username, password)

	if !mailService.MailService(data.Email, emailBody, "Your Account Credentials") {
		txn.Rollback()
		return fmt.Errorf("failed to send credentials email")
	}
	// ✅ Commit Transaction
	return txn.Commit().Error
}

func GetAllEmployeesService(db *gorm.DB) ([]model.EmployeeResponse, error) {
	var employees []model.EmployeeResponse

	query := `
		SELECT 
			u.*, 
			a."refUACUsername" AS username, 
			c."refUCDEmail" AS email, 
			c."refUCDMobile" AS mobile,
			c."refUCDDoorNo" AS doorNo, 
			c."refUCDStreet" AS street,
			c."refUCDCity" AS city,
			c."refUCDState" AS state
		FROM "Users" u
		LEFT JOIN "refUserAuthCred" a ON u."refUserId" = a."refUserId"
		LEFT JOIN "refUserCommunicationDetails" c ON u."refUserId" = c."refUserId"
		WHERE u."isDelete" = false
		ORDER BY u."refUserId" DESC;
	`

	if err := db.Raw(query).Scan(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch employees: %w", err)
	}

	return employees, nil
}

func GetEmployeeByIDService(db *gorm.DB, id string) (*model.EmployeeResponse, error) {
	var employee model.EmployeeResponse

	query := `
		SELECT 
			u.*, 
			a."refUACUsername" AS username, 
			c."refUCDEmail" AS email, 
			c."refUCDMobile" AS mobile,
			c."refUCDDoorNo" AS doorNo, 
			c."refUCDStreet" AS street,
			c."refUCDCity" AS city,
			c."refUCDState" AS state
		FROM "Users" u
		LEFT JOIN "refUserAuthCred" a ON u."refUserId" = a."refUserId"
		LEFT JOIN "refUserCommunicationDetails" c ON u."refUserId" = c."refUserId"
		WHERE u."isDelete" = false AND u."refUserId" = ?
		ORDER BY u."refUserId" ASC;
	`

	if err := db.Raw(query, id).Scan(&employee).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch employee: %w", err)
	}

	return &employee, nil
}

func UpdateEmployeeService(db *gorm.DB, id string, data *model.EmployeePayload) error {
	txn := db.Begin()
	if txn.Error != nil {
		return txn.Error
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	updatedBy := "Admin"

	// Step 1: Update Users
	userUpdate := map[string]interface{}{
		"refUserFName":       data.FirstName,
		"refUserLName":       data.LastName,
		"refUserDesignation": data.Designation,
		"refUserStatus":      map[bool]string{true: "Active", false: "In Active"}[data.RefUserStatus],
		"refUserBranchId":    data.RefUserBranchId,
		"updatedAt":          timestamp,
		"updatedBy":          updatedBy,
	}
	if err := txn.Table(`"Users"`).Where(`"refUserId" = ?`, id).Updates(userUpdate).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Step 2: Update Auth (optional - only username update allowed)
	if data.Username != "" {
		if err := txn.Table(`"refUserAuthCred"`).
			Where(`"refUserId" = ?`, id).
			Update("refUACUsername", data.Username).Error; err != nil {
			txn.Rollback()
			return fmt.Errorf("failed to update username: %w", err)
		}
	}

	// Step 3: Update Communication
	commUpdate := map[string]interface{}{
		"refUCDEmail":  data.Email,
		"refUCDMobile": data.Mobile,
		"refUCDDoorNo": data.DoorNumber,
		"refUCDStreet": data.StreetName,
		"refUCDCity":   data.City,
		"refUCDState":  data.State,
		"updatedAt":    timestamp,
		"updatedBy":    updatedBy,
	}
	if err := txn.Table(`"refUserCommunicationDetails"`).
		Where(`"refUserId" = ?`, id).Updates(commUpdate).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to update communication details: %w", err)
	}

	return txn.Commit().Error
}

func SoftDeleteEmployeeService(db *gorm.DB, id string) error {
	return db.Table(`"Users"`).Where(`"refUserId" = ?`, id).Update("isDelete", true).Error
}

func GetEmployeeService(db *gorm.DB, id string) (*model.EmployeeResponse, error) {
	var employee model.EmployeeResponse

	query := `
		SELECT 
			u.*, 
			a."refUACUsername" AS username, 
			c."refUCDEmail" AS email, 
			c."refUCDMobile" AS mobile,
			c."refUCDDoorNo" AS doorNo, 
			c."refUCDStreet" AS street,
			c."refUCDCity" AS city,
			c."refUCDState" AS state,
		r."refRTName" as role,
		b."refBranchName" as branch
			FROM public."Users" u
			LEFT JOIN public."refUserAuthCred" a ON u."refUserId" = a."refUserId"
			LEFT JOIN public."refUserCommunicationDetails" c ON u."refUserId" = c."refUserId"
		LEFT JOIN public."RoleType" r ON u."refRTId" = r."refRTId"
		LEFT JOIN public."Branches" b ON u."refUserBranchId" = b."refBranchId"
			WHERE u."isDelete" = false AND u."refUserId" = ?
			ORDER BY u."refUserId" ASC;
	`

	if err := db.Raw(query, id).Scan(&employee).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch employee: %w", err)
	}

	return &employee, nil
}

func UpdateProfileService(db *gorm.DB, id string, data *model.ProfilePayload) error {
	txn := db.Begin()
	if txn.Error != nil {
		return txn.Error
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	updatedBy := "Admin"

	// Step 1: Update Users
	userUpdate := map[string]interface{}{
		"refUserFName":       data.FirstName,
		"refUserLName":       data.LastName,
		"refUserDesignation": data.Designation,
		// "refUserStatus":      map[bool]string{true: "Active", false: "In Active"}[data.RefUserStatus],
		"refUserBranchId": data.RefUserBranchId,
		"updatedAt":       timestamp,
		"updatedBy":       updatedBy,
	}
	if err := txn.Table(`"Users"`).Where(`"refUserId" = ?`, id).Updates(userUpdate).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Step 2: Update Auth (optional - only username update allowed)
	if data.Username != "" {
		if err := txn.Table(`"refUserAuthCred"`).
			Where(`"refUserId" = ?`, id).
			Update("refUACUsername", data.Username).Error; err != nil {
			txn.Rollback()
			return fmt.Errorf("failed to update username: %w", err)
		}
	}

	// Step 3: Update Communication
	commUpdate := map[string]interface{}{
		"refUCDEmail":  data.Email,
		"refUCDMobile": data.Mobile,
		"refUCDDoorNo": data.DoorNumber,
		"refUCDStreet": data.StreetName,
		"refUCDCity":   data.City,
		"refUCDState":  data.State,
		"updatedAt":    timestamp,
		"updatedBy":    updatedBy,
	}
	if err := txn.Table(`"refUserCommunicationDetails"`).
		Where(`"refUserId" = ?`, id).Updates(commUpdate).Error; err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to update communication details: %w", err)
	}

	return txn.Commit().Error
}

func FetchSettingsOverview(db *gorm.DB) (model.SettingsOverview, error) {
	var overview model.SettingsOverview

	// Cards query
	cardsQuery := `
	SELECT
	  (SELECT COUNT(*) FROM public."Categories" WHERE "isDelete" = 'false' AND DATE_TRUNC('month', "createdAt"::Date) = DATE_TRUNC('month', CURRENT_DATE)) AS "Categories",
	  (SELECT COUNT(*) FROM public."Branches" WHERE "isDelete" = 'false' AND DATE_TRUNC('month', "createdAt"::date) = DATE_TRUNC('month', CURRENT_DATE)) AS "Branches",
	  (SELECT COUNT(*) FROM public."Supplier" WHERE "isDelete" = 'false' AND DATE_TRUNC('month', "createdAt"::date) = DATE_TRUNC('month', CURRENT_DATE)) AS "Supplier",
	  (SELECT COUNT(*) FROM public."Attributes" WHERE "isDelete" = 'false' AND DATE_TRUNC('month', "createdAt"::date) = DATE_TRUNC('month', CURRENT_DATE)) AS "Attributes";
	`

	if err := db.Raw(cardsQuery).Scan(&overview.Cards).Error; err != nil {
		return overview, fmt.Errorf("error fetching cards: %v", err)
	}

	// Latest suppliers
	if err := db.Raw(`SELECT * FROM public."Supplier" WHERE "isDelete" IS false ORDER BY "createdAt" DESC LIMIT 5`).Scan(&overview.LatestSuppliers).Error; err != nil {
		return overview, fmt.Errorf("error fetching latest suppliers: %v", err)
	}

	// Latest categories
	if err := db.Raw(`SELECT * FROM public."Categories" WHERE "isDelete" IS false ORDER BY "createdAt" DESC LIMIT 5`).Scan(&overview.LatestCategories).Error; err != nil {
		return overview, fmt.Errorf("error fetching latest categories: %v", err)
	}

	// Monthly category & subcategory counts
	monthlyQuery := `
	WITH category_counts AS (
		SELECT TO_CHAR("createdAt"::date, 'MM-YYYY') AS month, COUNT(*) AS "Categories"
		FROM public."Categories"
		WHERE "isDelete" = false
		GROUP BY TO_CHAR("createdAt"::date, 'MM-YYYY')
	),
	subcategory_counts AS (
		SELECT TO_CHAR("createdAt"::date, 'MM-YYYY') AS month, COUNT(*) AS "SubCategories"
		FROM public."SubCategories"
		WHERE "isDelete" = false
		GROUP BY TO_CHAR("createdAt"::date, 'MM-YYYY')
	)
	SELECT COALESCE(c.month, s.month) AS month,
	       COALESCE(c."Categories", 0) AS "Categories",
	       COALESCE(s."SubCategories", 0) AS "SubCategories"
	FROM category_counts c
	FULL OUTER JOIN subcategory_counts s ON c.month = s.month
	ORDER BY month;
	`

	if err := db.Raw(monthlyQuery).Scan(&overview.MonthlyCounts).Error; err != nil {
		return overview, fmt.Errorf("error fetching monthly counts: %v", err)
	}

	return overview, nil
}
