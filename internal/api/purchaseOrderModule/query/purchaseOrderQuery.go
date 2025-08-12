package purchaseOrderQuery

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GeneratePONumber(db *gorm.DB) (string, error) {
	var count int64
	err := db.Table(`"purchaseOrder"."CreatePurchaseOrder"`).Count(&count).Error
	if err != nil {
		return "", err
	}

	monthYear := time.Now().Format("0106")

	return fmt.Sprintf("PO-%s-%04d", monthYear, count+1), nil
}
