package roleType

import (
	"fmt"

	"gorm.io/gorm"
)

func GetRoleTypeNameByID(db *gorm.DB, refRTId int) (string, error) {
	var role struct {
		Name string `gorm:"column:refRTName"`
	}

	fmt.Println("refRTId", refRTId)
	err := db.Debug().
		Table(`"RoleType"`).
		Where(`"refRTId" = ?`, refRTId).
		Order(`"RoleType"."refRTName"`).
		First(&role).Error
	if err != nil {
		return "", fmt.Errorf("error fetching role name: %v", err)
	}

	fmt.Println("Role Name:", role.Name)
	return role.Name, nil
}
