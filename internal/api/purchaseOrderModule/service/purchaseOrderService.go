package purchaseOrderService

import (
	"fmt"

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

	// 2. Insert each Product with order.PurchaseOrderID
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
	}

	return nil
}

// purchaseOrderService/purchaseOrderService.go

func GetAllPurchaseOrdersService(db *gorm.DB) ([]purchaseOrderModel.CreatePORequest, error) {
	log := logger.InitLogger()

	log.Println("INFO: GetAllPurchaseOrdersService started")
	var orders []purchaseOrderModel.CreatePORequest

	type OrderRow struct {
		PurchaseOrderID int    `gorm:"column:purchaseOrderId"`
		PONumber        string `gorm:"column:poNumber"`
		SupplierID      int    `gorm:"column:supplierId"`
		BranchID        int    `gorm:"column:branchId"`
		Status          int    `gorm:"column:status"`
		ExpectedDate    string `gorm:"column:expectedDate"`
		ModeOfTransport string `gorm:"column:modeOfTransport"`
		SubTotal        string `gorm:"column:subTotal"`
		DiscountOverall string `gorm:"column:discountOverall"`
		PayAmount       string `gorm:"column:payAmount"`
		IsTaxApplied    bool   `gorm:"column:isTaxApplied"`
		TaxPercentage   string `gorm:"column:taxPercentage"`
		TaxedAmount     string `gorm:"column:taxedAmount"`
		TotalAmount     string `gorm:"column:totalAmount"`
		TotalPaid       string `gorm:"column:totalPaid"`
		PaymentPending  string `gorm:"column:paymentPending"`
		CreatedAt       string `gorm:"column:createdAt"`
		CreatedBy       string `gorm:"column:createdBy"`
		UpdatedAt       string `gorm:"column:updatedAt"`
		UpdatedBy       string `gorm:"column:updatedBy"`
		IsDelete        bool   `gorm:"column:isDelete"`

		// Supplier
		SupplierName          string `gorm:"column:supplierName"`
		SupplierCompanyName   string `gorm:"column:supplierCompanyName"`
		SupplierGSTNumber     string `gorm:"column:supplierGSTNumber"`
		SupplierEmail         string `gorm:"column:supplierEmail"`
		SupplierContactNumber string `gorm:"column:supplierContactNumber"`
		SupplierPaymentTerms  string `gorm:"column:supplierPaymentTerms"`
		SupplierAddress       string `gorm:"column:supplierAddress"`

		// Branch
		BranchName    string `gorm:"column:branchName"`
		BranchEmail   string `gorm:"column:branchEmail"`
		BranchAddress string `gorm:"column:branchAddress"`
	}

	var orderRows []OrderRow

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
			ProductDetails: products,
		}

		log.Printf("DEBUG: First row sample: %+v\n", orderRows[0])

		orders = append(orders, order)

		log.Info("\n\nINFO: Appended order for PurchaseOrderID = %d\n", row.PurchaseOrderID)
		fmt.Printf("DEBUG: order = %+v\n", order)
	}

	log.Println("INFO: GetAllPurchaseOrdersService completed successfully")
	fmt.Printf("DEBUG: Final orders list = %+v\n", orders)

	return orders, nil
}
