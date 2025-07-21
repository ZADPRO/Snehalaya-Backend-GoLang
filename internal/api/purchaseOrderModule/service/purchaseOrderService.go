package purchaseOrderService

import (
	"fmt"

	purchaseOrderModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/model"
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
	var orders []purchaseOrderModel.CreatePORequest

	type OrderRow struct {
		PurchaseOrderID int `gorm:"column:purchaseOrderId"`
		purchaseOrderModel.TotalSummary
		purchaseOrderModel.SupplierDetails
		purchaseOrderModel.BranchDetails
	}

	var orderRows []OrderRow
	if err := db.Table(`"purchaseOrder"."CreatePurchaseOrder" AS po`).
		Select(`po."purchaseOrderId", po.*, 
		s."supplierName", s."supplierCompanyName", s."supplierGSTNumber", s."supplierEmail", s."supplierContactNumber",
		s."supplierDoorNumber" || ', ' || s."supplierStreet" || ', ' || s."supplierCity" || ', ' || s."supplierState" || ', ' || s."supplierCountry" || ', PIN: ' || s."supplierPaymentTerms" AS supplierAddress,
		s."supplierPaymentTerms",
		b."refBranchName" AS branchName, b."refEmail" AS branchEmail, b."refLocation" AS branchAddress`).
		Joins(`LEFT JOIN "public"."Supplier" s ON po."supplierId" = s."supplierId"`).
		Joins(`LEFT JOIN "public"."Branches" b ON po."branchId" = b."refBranchId"`).
		Where(`po."isDelete" = 'false'`).
		Scan(&orderRows).Error; err != nil {
		return nil, err
	}

	for _, row := range orderRows {
		var products []purchaseOrderModel.ProductDetails
		if err := db.Table(`"purchaseOrder"."PurchaseOrderItemsInitial"`).
			Where(`"purchaseOrderId" = ? AND "isDelete" = ?`, row.PurchaseOrderID, false).
			Find(&products).Error; err != nil {
			return nil, err
		}

		order := purchaseOrderModel.CreatePORequest{
			SupplierDetails: row.SupplierDetails,
			BranchDetails:   row.BranchDetails,
			ProductDetails:  products,
			TotalSummary:    row.TotalSummary,
		}
		orders = append(orders, order)
	}

	return orders, nil
}
