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

	log.Infof("💾 Creating PO: %+v", po)
	if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).Create(&po).Error; err != nil {
		log.Error("❌ Failed to insert PO: " + err.Error())
		return err
	}
	log.Infof("✅ PO created successfully with ID: %d", po.PurchaseOrderId)

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

		for i := 0; i < prod.RejectedQty; i++ {
			instance := poModuleModel.PurchaseOrderProductInstances{
				PoProductID:        fmt.Sprintf("%v", poProduct.PoProductId),
				SerialNo:           fmt.Sprintf("%v", instanceNo),
				CategoryID:         prod.CategoryId,
				ProductDescription: prod.ProductName,
				UnitPrice:          prod.UnitPrice,
				Status:             "Rejected",
				CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:          roleName,
			}
			log.Infof("💾 Creating Product Instance (Rejected): %+v", instance)
			if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProductInstances"`).Create(&instance).Error; err != nil {
				log.Error("❌ Failed to insert Product Instance: " + err.Error())
				return err
			}
			instanceNo++
		}

		log.Infof("✅ Finished processing product #%d", idx+1)
	}

	log.Info("✅ PO, Products, and Product Instances created successfully")
	return nil
}
