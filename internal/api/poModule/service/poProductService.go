package poService

import (
	"encoding/json"
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

	// Step 1️⃣ - Check or Create PO
	var existingPO poModuleModel.PurchaseOrdersProducts
	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
		Where(`purchase_order_id = ? AND "invoiceNumber" = ?`, poPayload.PoId, poPayload.PoInvoiceNumber).
		First(&existingPO).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Info("⚠️ No existing PO found, creating new one.")

			newPO := poModuleModel.PurchaseOrdersProducts{
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

			if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).Create(&newPO).Error; err != nil {
				log.Error("❌ Failed to insert new PO: " + err.Error())
				return err
			}
			existingPO = newPO
			log.Infof("✅ New PO created with ID: %d", newPO.PurchaseOrderId)
		} else {
			log.Error("❌ Failed to fetch PO: " + err.Error())
			return err
		}
	} else {
		// Update existing PO
		log.Infof("🔄 Updating existing PO ID: %d", existingPO.PurchaseOrderId)
		existingPO.TotalAmount = poPayload.TotalAmount
		existingPO.InvoiceStatus = true
		existingPO.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		existingPO.UpdatedBy = roleName

		if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
			Where(`purchase_order_id = ?`, existingPO.PurchaseOrderId).
			Updates(&existingPO).Error; err != nil {
			log.Error("❌ Failed to update PO: " + err.Error())
			return err
		}
		log.Infof("✅ PO updated successfully with ID: %d", existingPO.PurchaseOrderId)
	}

	// Step 2️⃣ - Insert or Update PO Products
	for idx, prod := range poPayload.Products {
		log.Infof("💡 Processing product #%d: %+v", idx+1, prod)
		status := "Accepted"
		if prod.RejectedQty > 0 {
			status = "Partial"
		}

		var existingProduct poModuleModel.PurchaseOrderProducts
		err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts"`).
			Where(`purchase_order_id = ? AND category_id = ? AND description = ?`,
				existingPO.PurchaseOrderId, prod.CategoryId, prod.ProductName).
			First(&existingProduct).Error

		if err == nil {
			// 🔄 Update existing product
			existingProduct.UnitPrice = prod.UnitPrice
			existingProduct.Quantity = fmt.Sprintf("%v", prod.OrderedQty)
			existingProduct.AcceptedQuantity = fmt.Sprintf("%v", prod.ReceivedQty)
			existingProduct.RejectedQuantity = fmt.Sprintf("%v", prod.RejectedQty)
			existingProduct.Status = status
			existingProduct.AcceptedTotal = fmt.Sprintf("%v", prod.ReceivedQty)
			existingProduct.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			existingProduct.UpdatedBy = roleName

			if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts"`).
				Where(`po_product_id = ?`, existingProduct.PoProductId).
				Updates(&existingProduct).Error; err != nil {
				log.Error("❌ Failed to update PO Product: " + err.Error())
				return err
			}
			log.Infof("✅ Updated existing product: %s", prod.ProductName)
		} else if err == gorm.ErrRecordNotFound {
			// ➕ Insert new product
			newProduct := poModuleModel.PurchaseOrderProducts{
				PurchaseOrderID:  existingPO.PurchaseOrderId,
				CategoryID:       prod.CategoryId,
				Description:      prod.ProductName,
				UnitPrice:        prod.UnitPrice,
				Quantity:         fmt.Sprintf("%v", prod.OrderedQty),
				AcceptedQuantity: fmt.Sprintf("%v", prod.ReceivedQty),
				RejectedQuantity: fmt.Sprintf("%v", prod.RejectedQty),
				Status:           status,
				AcceptedTotal:    fmt.Sprintf("%v", prod.ReceivedQty),
				CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:        roleName,
			}
			if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts"`).Create(&newProduct).Error; err != nil {
				log.Error("❌ Failed to insert new PO Product: " + err.Error())
				return err
			}
			existingProduct = newProduct
			log.Infof("✅ Created new product: %s", prod.ProductName)
		} else {
			log.Error("❌ Error checking existing product: " + err.Error())
			return err
		}

		// Step 3️⃣ - Manage Product Instances
		db.Table(`"purchaseOrderMgmt"."PurchaseOrderProductInstances"`).
			Where(`po_product_id = ?`, existingProduct.PoProductId).
			Delete(nil)

		for i := 1; i <= prod.ReceivedQty; i++ {
			instance := poModuleModel.PurchaseOrderProductInstances{
				PoProductID:        fmt.Sprintf("%v", existingProduct.PoProductId),
				SerialNo:           fmt.Sprintf("%v", i),
				CategoryID:         prod.CategoryId,
				ProductDescription: prod.ProductName,
				UnitPrice:          prod.UnitPrice,
				Status:             "Accepted",
				CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:          roleName,
			}
			if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProductInstances"`).Create(&instance).Error; err != nil {
				log.Error("❌ Failed to insert product instance: " + err.Error())
				return err
			}
		}

		// Step 4️⃣ - Store Rejected Products
		if prod.RejectedQty > 0 {
			rejected := poModuleModel.RejectedProducts{
				PoProductID:        existingProduct.PoProductId,
				CategoryID:         prod.CategoryId,
				ProductDescription: prod.ProductName,
				UnitPrice:          prod.UnitPrice,
				RejectedQty:        fmt.Sprintf("%v", prod.RejectedQty),
				Reason:             "",
				CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
				CreatedBy:          roleName,
			}
			if err := db.Table(`"purchaseOrderMgmt"."RejectedProducts"`).Create(&rejected).Error; err != nil {
				log.Error("❌ Failed to insert rejected product: " + err.Error())
				return err
			}
			log.Infof("🚫 Rejected product recorded: %s (%d pcs)", prod.ProductName, prod.RejectedQty)
		}
	}

	log.Info("✅ PO, Products, and Product Instances processed successfully")
	return nil
}

func GetAcceptedPurchaseOrdersService(db *gorm.DB) ([]poModuleModel.AcceptedPOResponse, error) {
	log := logger.InitLogger()
	log.Info("🧾 Fetching accepted purchase orders...")

	type rawPO struct {
		PurchaseOrderID  int     `json:"purchase_order_id"`
		InvoiceNumber    string  `json:"invoice_number"`
		BranchID         int     `json:"branch_id"`
		SupplierID       int     `json:"supplier_id"`
		TotalAmount      string  `json:"total_amount"`
		CreatedAt        string  `json:"created_at"`
		AcceptedProducts *string `json:"accepted_products"` // raw JSON string
	}

	var rawResults []rawPO
	query := `
		SELECT 
			po.purchase_order_id,
			po."invoiceNumber" AS invoice_number,
			po.branch_id,
			po.supplier_id,
			po.total_amount,
			po."createdAt" AS created_at,
			COALESCE(
				json_agg(
					json_build_object(
						'po_product_id', p.po_product_id,
						'category_id', p.category_id,
						'product_description', p.description,
						'unit_price', p.unit_price,
						'ordered_quantity', p.quantity,
						'ordered_total', p.total,
						'accepted_quantity', p.accepted_quantity,
						'accepted_total', p.accepted_total,
						'status', p.status,
						'updated_at', p."updatedAt",
						'updated_by', p."updatedBy"
					)
				) FILTER (WHERE p.accepted_quantity::int > 0), '[]'
			) AS accepted_products
		FROM "purchaseOrderMgmt"."PurchaseOrders" po
		LEFT JOIN "purchaseOrderMgmt"."PurchaseOrderProducts" p
			ON po.purchase_order_id = p.purchase_order_id
		GROUP BY po.purchase_order_id, po."invoiceNumber", po.branch_id, po.supplier_id, po.total_amount, po."createdAt"
		ORDER BY po.purchase_order_id DESC;
	`

	if err := db.Raw(query).Scan(&rawResults).Error; err != nil {
		log.Error("❌ Query execution failed: " + err.Error())
		return nil, err
	}

	var finalResults []poModuleModel.AcceptedPOResponse
	for _, r := range rawResults {
		var products []poModuleModel.AcceptedProduct
		if r.AcceptedProducts != nil && *r.AcceptedProducts != "" {
			if err := json.Unmarshal([]byte(*r.AcceptedProducts), &products); err != nil {
				log.Warnf("⚠️ Failed to parse products for PO ID %d: %v", r.PurchaseOrderID, err)
			}
		}
		finalResults = append(finalResults, poModuleModel.AcceptedPOResponse{
			PurchaseOrderID:  r.PurchaseOrderID,
			InvoiceNumber:    r.InvoiceNumber,
			BranchID:         r.BranchID,
			SupplierID:       r.SupplierID,
			TotalAmount:      r.TotalAmount,
			CreatedAt:        r.CreatedAt,
			AcceptedProducts: products,
		})
	}

	log.Infof("✅ Retrieved %d accepted purchase orders", len(finalResults))
	return finalResults, nil
}
