package purchaseOrderService

import (
	"fmt"
	"strconv"
	"time"

	purchaseOrderModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreatePurchaseOrderService(db *gorm.DB, payload *purchaseOrderModel.CreatePORequest, createdBy string) error {
	// 1. Insert into CreatePurchaseOrder
	order := purchaseOrderModel.CreatePurchaseOrder{
		PONumber:        payload.TotalSummary.PONumber,
		SupplierID:      payload.TotalSummary.SupplierID,
		BranchID:        payload.TotalSummary.BranchID,
		Status:          payload.TotalSummary.Status,
		ExpectedDate:    payload.TotalSummary.ExpectedDate,
		ModeOfTransport: payload.TotalSummary.ModeOfTransport,
		SubTotal:        payload.TotalSummary.SubTotal,
		DiscountOverall: payload.TotalSummary.DiscountOverall,
		PayAmount:       payload.TotalSummary.PayAmount,
		IsTaxApplied:    payload.TotalSummary.IsTaxApplied,
		TaxPercentage:   payload.TotalSummary.TaxPercentage,
		TaxedAmount:     payload.TotalSummary.TaxedAmount,
		TotalAmount:     payload.TotalSummary.TotalAmount,
		TotalPaid:       payload.TotalSummary.TotalPaid,
		PaymentPending:  payload.TotalSummary.PaymentPending,
		CreatedAt:       payload.TotalSummary.CreatedAt,
		CreatedBy:       createdBy,
		UpdatedAt:       payload.TotalSummary.UpdatedAt,
		UpdatedBy:       payload.TotalSummary.UpdatedBy,
		IsDelete:        fmt.Sprintf("%v", payload.TotalSummary.IsDelete),
	}

	if err := db.Create(&order).Error; err != nil {
		return err
	}

	// 2. Insert Product Items and Dummy Acceptance Rows
	for _, item := range payload.ProductDetails {
		dbItem := purchaseOrderModel.PurchaseOrderItem{
			PurchaseOrderID:  order.PurchaseOrderID,
			ProductName:      item.ProductName,
			RefCategoryID:    item.RefCategoryID,
			RefSubCategoryID: item.RefSubCategoryID,
			HSNCode:          item.HSNCode,
			PurchaseQuantity: item.PurchaseQuantity,
			PurchasePrice:    item.PurchasePrice,
			DiscountPrice:    item.DiscountPrice,
			DiscountAmount:   item.DiscountAmount,
			TotalAmount:      item.TotalAmount,
			IsReceived:       item.IsReceived,
			AcceptanceStatus: item.AcceptanceStatus,
			CreatedAt:        item.CreatedAt,
			CreatedBy:        item.CreatedBy,
			UpdatedAt:        item.UpdatedAt,
			UpdatedBy:        item.UpdatedBy,
			IsDelete:         item.IsDelete,
		}

		if err := db.Create(&dbItem).Error; err != nil {
			return err
		}

		// 3. Insert N dummy product entries
		qty, err := strconv.Atoi(item.PurchaseQuantity)
		if err != nil {
			return fmt.Errorf("invalid purchaseQuantity for product %s: %v", item.ProductName, err)
		}

		for i := 0; i < qty; i++ {
			dummy := purchaseOrderModel.ProductsDummyAcceptance{
				PurchaseOrderID:  order.PurchaseOrderID,
				ProductName:      item.ProductName,
				RefCategoryID:    item.RefCategoryID,
				RefSubCategoryID: item.RefSubCategoryID,
				HSNCode:          item.HSNCode,
				DummySKU:         fmt.Sprintf("%s-%d", item.ProductName, i+1), // Or generate your own SKU logic
				Price:            item.PurchasePrice,
				DiscountAmount:   item.DiscountAmount,
				DiscountPercent:  item.DiscountPrice, // Assuming you map accordingly
				IsReceived:       "false",
				AcceptanceStatus: "Pending",
				CreatedAt:        item.CreatedAt,
				CreatedBy:        createdBy,
				UpdatedAt:        item.UpdatedAt,
				UpdatedBy:        item.UpdatedBy,
				IsDelete:         "false",
			}

			if err := db.Create(&dummy).Error; err != nil {
				return fmt.Errorf("failed to create dummy product for %s: %v", item.ProductName, err)
			}
		}
	}

	return nil
}

// purchaseOrderService/purchaseOrderService.go

func GetAllPurchaseOrdersService(db *gorm.DB) ([]purchaseOrderModel.CreatePORequest, error) {
	log := logger.InitLogger()

	log.Println("INFO: GetAllPurchaseOrdersService started")
	var orders []purchaseOrderModel.CreatePORequest
	var orderRows []purchaseOrderModel.OrderRow

	query := `
		SELECT
			po."purchaseOrderId",
			po."poNumber", po."supplierId", po."branchId", po."status",
			po."expectedDate", po."modeOfTransport", po."subTotal", po."discountOverall", po."payAmount",
			po."isTaxApplied", po."taxPercentage", po."taxedAmount", po."totalAmount",
			po."totalPaid", po."paymentPending",
			po."createdAt", po."createdBy", po."updatedAt", po."updatedBy", po."isDelete",

			s."supplierId" AS "supplierId",
			s."supplierName", s."supplierCompanyName", s."supplierGSTNumber",
			s."supplierEmail", s."supplierContactNumber",
			s."supplierPaymentTerms",
			CONCAT_WS(', ',
				s."supplierDoorNumber",
				s."supplierStreet",
				s."supplierCity",
				s."supplierState",
				s."supplierCountry"
			) AS "supplierAddress",

			b."refBranchId" AS "branchId",
			b."refBranchName" AS "branchName",
			b."refEmail" AS "branchEmail",
			b."refLocation" AS "branchAddress"

		FROM "purchaseOrder"."CreatePurchaseOrder" po
		LEFT JOIN "public"."Supplier" s ON po."supplierId" = s."supplierId" AND s."isDelete" = false
		LEFT JOIN "public"."Branches" b ON po."branchId" = b."refBranchId" AND b."isDelete" = false
		WHERE po."isDelete" = 'false';
	`

	if err := db.Raw(query).Scan(&orderRows).Error; err != nil {
		log.Println("ERROR: Failed to fetch purchase orders:", err)
		return nil, err
	}

	log.Println("\n\nINFO: Purchase order rows fetched successfully")
	fmt.Printf("DEBUG: orderRows = %+v\n", orderRows)

	for _, row := range orderRows {
		var products []purchaseOrderModel.ProductDetails

		if err := db.Raw(`
			SELECT *
			FROM "purchaseOrder"."PurchaseOrderItemsInitial"
			WHERE "purchaseOrderId" = ? AND "isDelete" = false
		`, row.PurchaseOrderID).Scan(&products).Error; err != nil {
			log.Info("\n\nERROR: Failed to fetch products for PurchaseOrderID =", row.PurchaseOrderID, ":", err)
			return nil, err
		}

		log.Printf("\n\nINFO: Fetched %d products for PurchaseOrderID = %d\n", len(products), row.PurchaseOrderID)
		fmt.Printf("DEBUG: products for PurchaseOrderID %d = %+v\n", row.PurchaseOrderID, products)

		order := purchaseOrderModel.CreatePORequest{
			SupplierDetails: purchaseOrderModel.SupplierDetails{
				SupplierID:            row.SupplierID,
				SupplierName:          row.SupplierName,
				SupplierCompanyName:   row.SupplierCompanyName,
				SupplierGSTNumber:     row.SupplierGSTNumber,
				SupplierAddress:       row.SupplierAddress,
				SupplierPaymentTerms:  row.SupplierPaymentTerms,
				SupplierEmail:         row.SupplierEmail,
				SupplierContactNumber: row.SupplierContactNumber,
			},
			BranchDetails: purchaseOrderModel.BranchDetails{
				BranchID:      row.BranchID,
				BranchName:    row.BranchName,
				BranchEmail:   row.BranchEmail,
				BranchAddress: row.BranchAddress,
			},
			TotalSummary: purchaseOrderModel.TotalSummary{
				PONumber:        row.PONumber,
				SupplierID:      row.SupplierID,
				BranchID:        row.BranchID,
				Status:          row.Status,
				ExpectedDate:    row.ExpectedDate,
				ModeOfTransport: row.ModeOfTransport,
				SubTotal:        row.SubTotal,
				DiscountOverall: row.DiscountOverall,
				PayAmount:       row.PayAmount,
				IsTaxApplied:    row.IsTaxApplied,
				TaxPercentage:   row.TaxPercentage,
				TaxedAmount:     row.TaxedAmount,
				TotalAmount:     row.TotalAmount,
				TotalPaid:       row.TotalPaid,
				PaymentPending:  row.PaymentPending,
				CreatedAt:       row.CreatedAt,
				CreatedBy:       row.CreatedBy,
				UpdatedAt:       row.UpdatedAt,
				UpdatedBy:       row.UpdatedBy,
				IsDelete:        row.IsDelete,
			},
			ProductDetails:  products,
			PurchaseOrderID: row.PurchaseOrderID,
		}

		log.Printf("\n\nDEBUG: First row sample: %+v\n", orderRows[0])

		orders = append(orders, order)

		fmt.Printf("\n\nDEBUG: order = %+v\n", order)
	}

	log.Println("INFO: GetAllPurchaseOrdersService completed successfully")
	fmt.Printf("DEBUG: Final orders list = %+v\n", orders)

	return orders, nil
}

func GetDummyProductsByPOIDService(db *gorm.DB, poID string) ([]purchaseOrderModel.ProductsDummyAcceptance, error) {
	var dummyProducts []purchaseOrderModel.ProductsDummyAcceptance

	if err := db.
		Where(`"purchaseOrderId" = ? AND "isDelete" = ?`, poID, "false").
		Find(&dummyProducts).Error; err != nil {
		return nil, err
	}

	return dummyProducts, nil
}

func UpdateDummyProductStatusService(db *gorm.DB, dummyProductId int, status interface{}, reason string) error {
	var product purchaseOrderModel.ProductsDummyAcceptance

	if err := db.First(&product, dummyProductId).Error; err != nil {
		return fmt.Errorf("dummy product not found: %v", err)
	}

	switch val := status.(type) {
	case bool:
		if val {
			// ✅ Accept
			if product.IsReceived != "true" {
				now := time.Now()
				month := fmt.Sprintf("%02d", now.Month())
				year := fmt.Sprintf("%02d", now.Year()%100)

				var count int64
				db.Model(&purchaseOrderModel.ProductsDummyAcceptance{}).
					Where("dummySKU LIKE ?", fmt.Sprintf("SS%s%s%%", month, year)).
					Count(&count)

				sku := fmt.Sprintf("SS%s%s%05d", month, year, count+1)

				product.DummySKU = sku
				product.IsReceived = "true"
				product.AcceptanceStatus = "Received"
			}
		} else {
			// ❌ Reject
			product.DummySKU = ""
			product.IsReceived = "false"
			product.AcceptanceStatus = reason
		}
	case string:
		if val == "undo" {
			// 🔄 Undo
			product.DummySKU = ""
			product.IsReceived = "false"
			product.AcceptanceStatus = "Pending"
		}
	default:
		return fmt.Errorf("invalid status type")
	}

	return db.Save(&product).Error
}

// BULK UPDATE - ACCEPT, REJECT, UNDO
func BulkUpdateDummyProducts(db *gorm.DB, ids []int, action string, reason string) error {
	var products []purchaseOrderModel.ProductsDummyAcceptance
	if err := db.Where(`"dummyProductsId" IN ?`, ids).Find(&products).Error; err != nil {
		return err
	}

	now := time.Now()
	month := fmt.Sprintf("%02d", now.Month())
	year := fmt.Sprintf("%02d", now.Year()%100)

	switch action {
	case "accept":
		// Get latest SKU count
		var count int64
		db.Model(&purchaseOrderModel.ProductsDummyAcceptance{}).
			Where("dummySKU LIKE ?", fmt.Sprintf("SS%s%s%%", month, year)).
			Count(&count)

		for i := range products {
			if products[i].IsReceived != "true" {
				count++
				products[i].DummySKU = fmt.Sprintf("SS%s%s%05d", month, year, count)
				products[i].IsReceived = "true"
				products[i].AcceptanceStatus = "Received"
			}
		}
	case "reject":
		for i := range products {
			products[i].DummySKU = ""
			products[i].IsReceived = "false"
			products[i].AcceptanceStatus = reason
		}
	case "undo":
		for i := range products {
			products[i].DummySKU = ""
			products[i].IsReceived = "false"
			products[i].AcceptanceStatus = "Pending"
		}
	}

	// Save all
	for _, product := range products {
		if err := db.Save(&product).Error; err != nil {
			return err
		}
	}
	return nil
}
