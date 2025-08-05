package roleType

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func ExtractIntFromInterface(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New("string cannot be converted to int")
		}
		return parsed, nil
	default:
		return 0, errors.New("unsupported type for conversion to int")
	}
}

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
