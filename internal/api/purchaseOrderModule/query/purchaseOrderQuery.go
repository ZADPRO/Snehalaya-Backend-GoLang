package purchaseOrderQuery

import (
	"fmt"

	"gorm.io/gorm"

)

func GeneratePONumber(db *gorm.DB) (string, error) {
	var count int64
	err := db.Table(`"purchaseOrder"."CreatePurchaseOrder"`).Count(&count).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PO-2025-%04d", count+1), nil
}
