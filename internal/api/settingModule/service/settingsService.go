package settingsService

import (
	"errors"
	"fmt"
	"time"

	transactionLogger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/model"
	becrypt "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Bcrypt"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	mailService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/MailService"
	"gorm.io/gorm"
)

// CATEGORIES SERVICE

func CreateCategoryService(db *gorm.DB, category *model.Category) error {
	log := logger.InitLogger()

	// Check for existing category with same name or code and isDelete = false
	var existing model.Category

	fmt.Println("category", category)

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
	err = db.Table("Categories").Create(category).Error
	if err == nil {
		_ = transactionLogger.LogTransaction(db, 1, "Admin", 2, "Category Created: "+category.CategoryName)
	}
	return err
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
	_ = transactionLogger.LogTransaction(db, 1, "Admin", 3, "Category Updated: "+category.CategoryName)

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

	err := db.Table(`"Branches"`).Where(`"isDelete" = false`).Find(&branches).Error
	if err != nil {
		log.Error("Failed to fetch branches: " + err.Error())
	}
	return branches
}

func UpdateBranchService(db *gorm.DB, branch *model.Branch) error {
	log := logger.InitLogger()
	log.Info("Updating Branch ID: ", branch.RefBranchId)

	var existing model.Branch
	err := db.Table(`"Branches"`).
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
		"updatedBy":     "Admin",
	}

	return db.Table(`"Branches"`).Where(`"refBranchId" = ?`, branch.RefBranchId).Updates(updateData).Error
}

func DeleteBranchService(db *gorm.DB, id string) error {
	log := logger.InitLogger()
	log.Info("Soft deleting Branch with ID: ", id)

	return db.Table(`"Branches"`).
		Where(`"refBranchId" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": "Admin",
		}).Error
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

func GetAllEmployeesService(db *gorm.DB) ([]model.User, error) {
	var users []model.User
	if err := db.Where(`"isDelete" = ?`, false).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetEmployeeByIDService(db *gorm.DB, id string) (*model.EmployeeResponse, error) {
	var user model.User
	if err := db.Table("Users").
		Where(`"refUserId" = ? AND "isDelete" = ?`, id, false).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	var comm model.UserCommunication
	if err := db.Table(`"refUserCommunicationDetails"`).
		Where(`"refUserId" = ?`, id).
		First(&comm).Error; err != nil {
		return nil, fmt.Errorf("communication details not found")
	}

	var auth model.UserAuth
	if err := db.Table(`"refUserAuthCred"`).
		Where(`"refUserId" = ?`, id).
		First(&auth).Error; err != nil {
		return nil, fmt.Errorf("auth details not found")
	}

	response := &model.EmployeeResponse{
		User:     user,
		Username: auth.RefUACUsername,
		Email:    comm.RefUCDEmail,
		Mobile:   comm.RefUCDMobile,
		DoorNo:   comm.RefUCDDoorNo,
		Street:   comm.RefUCDStreet,
		City:     comm.RefUCDCity,
		State:    comm.RefUCDState,
	}

	return response, nil
}

func UpdateEmployeeService(db *gorm.DB, id string, data *model.EmployeePayload) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Where(`"refUserId" = ? AND "isDelete" = ?`, id, false).First(&user).Error; err != nil {
			return fmt.Errorf("employee not found")
		}

		user.RefUserFName = data.FirstName
		user.RefUserLName = data.LastName
		user.RefUserDesignation = data.Designation
		user.RefRTId = data.RoleTypeId
		user.RefUserStatus = map[bool]string{true: "Active", false: "In Active"}[data.RefUserStatus]
		user.RefUserBranchId = data.RefUserBranchId
		user.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		user.UpdatedBy = "Admin"

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Also update communication
		if err := tx.Model(&model.UserCommunication{}).
			Where(`"refUserId" = ?`, id).
			Updates(map[string]interface{}{
				"refUCDEmail":  data.Email,
				"refUCDMobile": data.Mobile,
				"refUCDDoorNo": data.DoorNumber,
				"refUCDStreet": data.StreetName,
				"refUCDCity":   data.City,
				"refUCDState":  data.State,
				"updatedAt":    time.Now().Format("2006-01-02 15:04:05"),
				"updatedBy":    "Admin",
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func SoftDeleteEmployeeService(db *gorm.DB, id string) error {
	return db.Model(&model.User{}).
		Where(`"refUserId" = ?`, id).
		Update(`"isDelete"`, true).Error
}
