package poService

import (
	"fmt"
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/service"
	poModuleModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreatePurchaseOrderService(db *gorm.DB, poPayload *poModuleModel.PurchaseOrderPayload, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ CreatePurchaseOrderService invoked")

	// 1️⃣ Create Purchase Order
	po := poModuleModel.PurchaseOrder{
		SupplierID:    poPayload.Supplier.SupplierId,
		BranchID:      poPayload.Branch.RefBranchId,
		SubTotal:      fmt.Sprintf("%v", poPayload.Summary.SubTotal),
		TotalDiscount: fmt.Sprintf("%v", poPayload.Summary.TotalDiscount),
		TaxEnabled:    poPayload.Summary.TaxEnabled,
		TaxPercentage: fmt.Sprintf("%v", poPayload.Summary.TaxPercentage),
		TaxAmount:     fmt.Sprintf("%v", poPayload.Summary.TaxAmount),
		TotalAmount:   fmt.Sprintf("%v", poPayload.Summary.TotalAmount),
		CreditedDate:  poPayload.CreditedDate,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:     roleName,
		IsDelete:      false,
	}

	if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).Create(&po).Error; err != nil {
		log.Error("❌ Failed to create Purchase Order: " + err.Error())
		return err
	}

	log.Infof("✅ Purchase Order created with ID: %d", po.PurchaseOrderID)

	// 2️⃣ Insert Products
	for _, prod := range poPayload.Products {
		fmt.Printf("DEBUG prod: %+v\n", prod) // shows all fields with names
		fmt.Printf("UnitPrice type: %T, value: %v\n", prod.UnitPrice, prod.UnitPrice)

		product := poModuleModel.PurchaseOrderProduct{
			PurchaseOrderID: po.PurchaseOrderID,
			CategoryID:      prod.CategoryID,
			Description:     prod.Description,
			UnitPrice:       fmt.Sprintf("%v", prod.UnitPrice),
			Discount:        fmt.Sprintf("%v", prod.Discount),
			Quantity:        fmt.Sprintf("%v", prod.Quantity),
			Total:           fmt.Sprintf("%v", prod.Total),
			CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:       roleName,
		}

		if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts"`).Create(&product).Error; err != nil {
			log.Error("❌ Failed to insert product: " + err.Error())
			return err
		}
	}

	transErr := service.LogTransaction(db, 1, "Admin", 2, fmt.Sprintf("PO Created: %d", po.PurchaseOrderID))
	if transErr != nil {
		log.Error("Failed to log transaction : " + transErr.Error())

	} else {
		log.Info("Transaction Log saved Successfully \n\n")
	}

	log.Info("✅ Purchase Order and Products saved successfully")
	return nil
}

func GetAllPurchaseOrdersService(db *gorm.DB) []poModuleModel.PurchaseOrderResponse {
	log := logger.InitLogger()
	var pos []poModuleModel.PurchaseOrderResponse

	err := db.Table("PurchaseOrders AS po").
		Select(`po.purchase_order_id, po.supplier_id, po.branch_id, po.sub_total, po.total_discount, 
				po.tax_enabled, po.tax_percentage, po.tax_amount, po.total_amount, po.credited_date,
				po.createdAt, po.createdBy`).
		Where("po.isDelete = ?", false).
		Order("po.purchase_order_id DESC").
		Scan(&pos).Error

	if err != nil {
		log.Error("❌ Failed to fetch purchase orders: " + err.Error())
	}

	return pos
}

func UpdatePurchaseOrderService(db *gorm.DB, poPayload *poModuleModel.PurchaseOrderPayload, roleName string) error {
	log := logger.InitLogger()

	updateData := map[string]interface{}{
		"supplier_id":    poPayload.Supplier.SupplierId,
		"branch_id":      poPayload.Branch.RefBranchId,
		"sub_total":      fmt.Sprintf("%v", poPayload.Summary.SubTotal),
		"total_discount": fmt.Sprintf("%v", poPayload.Summary.TotalDiscount),
		"tax_enabled":    poPayload.Summary.TaxEnabled,
		"tax_percentage": fmt.Sprintf("%v", poPayload.Summary.TaxPercentage),
		"tax_amount":     fmt.Sprintf("%v", poPayload.Summary.TaxAmount),
		"total_amount":   fmt.Sprintf("%v", poPayload.Summary.TotalAmount),
		"credited_date":  poPayload.CreditedDate,
		"updatedAt":      time.Now().Format("2006-01-02 15:04:05"),
		"updatedBy":      roleName,
	}

	if err := db.Table("PurchaseOrders").
		Where("purchase_order_id = ?", poPayload.PurchaseOrderID).
		Updates(updateData).Error; err != nil {
		log.Error("❌ Failed to update PO: " + err.Error())
		return err
	}

	log.Info("✅ PO updated successfully")
	return nil
}
