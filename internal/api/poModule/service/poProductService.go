package poService

import (
	"fmt"
	"time"

	poModuleModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreatePurchaseOrderProductService(db *gorm.DB, poPayload *poModuleModel.PurchaseOrderProductPayload, roleName string) error {
	log := logger.InitLogger()
	log.Info("🛠️ CreatePurchaseOrderService invoked")
	log.Infof("📦 Received PO Payload: %+v", poPayload)

	po := poModuleModel.PurchaseOrdersProducts{
		SupplierID:    poPayload.SupplierId,
		BranchID:      poPayload.BranchId,
		TotalAmount:   poPayload.TotalAmount,
		CreditedDate:  time.Now().Format("2006-01-02 15:04:05"),
		InvoiceNumber: poPayload.PoInvoiceNumber,
		InvoiceStatus: true,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:     roleName,
	}

	// Step 1️⃣ - Check if PO exists
	var existingPO poModuleModel.PurchaseOrdersProducts
	if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
		Where(`purchase_order_id = ? AND "invoiceNumber" = ?`, poPayload.PoId, poPayload.PoInvoiceNumber).
		First(&existingPO).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			log.Infof("⚠️ No existing PO found, creating new one instead.")
			po := poModuleModel.PurchaseOrdersProducts{
				PurchaseOrderId: poPayload.PoId,
				SupplierID:      poPayload.SupplierId,
				BranchID:        poPayload.BranchId,
				TotalAmount:     poPayload.TotalAmount,
				CreditedDate:    time.Now().Format("2006-01-02 15:04:05"),
				InvoiceNumber:   poPayload.PoInvoiceNumber,
				InvoiceStatus:   true,
				CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:       roleName,
			}

			if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).Create(&po).Error; err != nil {
				log.Error("❌ Failed to insert PO: " + err.Error())
				return err
			}
			log.Infof("✅ PO created successfully with ID: %d", po.PurchaseOrderId)
			existingPO = po
		} else {
			log.Error("❌ Failed to fetch PO: " + err.Error())
			return err
		}

	} else {
		// Update existing PO
		log.Infof("🔄 Updating existing PO ID: %d, Invoice: %s", existingPO.PurchaseOrderId, existingPO.InvoiceNumber)
		existingPO.TotalAmount = poPayload.TotalAmount
		existingPO.InvoiceStatus = true
		existingPO.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		existingPO.UpdatedBy = roleName

		if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
			Where(`purchase_order_id = ?`, poPayload.PoId).
			Updates(&existingPO).Error; err != nil {
			log.Error("❌ Failed to update existing PO: " + err.Error())
			return err
		}

		log.Infof("✅ PO updated successfully with ID: %d", existingPO.PurchaseOrderId)
	}

	// Insert PO Products
	for idx, prod := range poPayload.Products {
		log.Infof("💡 Processing product #%d: %+v", idx+1, prod)
		totalAccepted := fmt.Sprintf("%v", prod.ReceivedQty)
		status := "Accepted"
		if prod.RejectedQty > 0 {
			status = "Partial"
		}

		poProduct := poModuleModel.PurchaseOrderProducts{
			PurchaseOrderID:  po.PurchaseOrderId,
			CategoryID:       prod.CategoryId,
			Description:      prod.ProductName,
			UnitPrice:        prod.UnitPrice,
			Quantity:         fmt.Sprintf("%v", prod.OrderedQty),
			AcceptedQuantity: fmt.Sprintf("%v", prod.ReceivedQty),
			RejectedQuantity: fmt.Sprintf("%v", prod.RejectedQty),
			Status:           status,
			AcceptedTotal:    totalAccepted,
			CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:        roleName,
		}

		log.Infof("💾 Creating PO Product: %+v", poProduct)
		if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts"`).Create(&poProduct).Error; err != nil {
			log.Error("❌ Failed to insert PO Product: " + err.Error())
			return err
		}
		log.Infof("✅ PO Product created successfully with ID: %d", poProduct.PoProductId)

		// Generate Product Instances
		instanceNo := 1
		for i := 0; i < prod.ReceivedQty; i++ {
			instance := poModuleModel.PurchaseOrderProductInstances{
				PoProductID:        fmt.Sprintf("%v", poProduct.PoProductId),
				SerialNo:           fmt.Sprintf("%v", instanceNo),
				CategoryID:         prod.CategoryId,
				ProductDescription: prod.ProductName,
				UnitPrice:          prod.UnitPrice,
				Status:             "Accepted",
				CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:          roleName,
			}
			log.Infof("💾 Creating Product Instance (Accepted): %+v", instance)
			if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProductInstances"`).Create(&instance).Error; err != nil {
				log.Error("❌ Failed to insert Product Instance: " + err.Error())
				return err
			}
			instanceNo++
		}

		if prod.RejectedQty > 0 {
			rejected := poModuleModel.RejectedProducts{
				PoProductID:        poProduct.PoProductId,
				CategoryID:         prod.CategoryId,
				ProductDescription: prod.ProductName,
				UnitPrice:          prod.UnitPrice,
				RejectedQty:        fmt.Sprintf("%v", prod.RejectedQty),
				Reason:             "", // Optional: can pass rejection reason
				CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:          roleName,
			}
			log.Infof("💾 Creating Rejected Product: %+v", rejected)
			if err := db.Table(`"purchaseOrderMgmt"."RejectedProducts"`).Create(&rejected).Error; err != nil {
				log.Error("❌ Failed to insert Rejected Product: " + err.Error())
				return err
			}
		}

		log.Infof("✅ Finished processing product #%d", idx+1)
	}

	log.Info("✅ PO, Products, and Product Instances created successfully")
	return nil
}
