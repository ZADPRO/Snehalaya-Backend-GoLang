package poService

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/service"
	poModuleModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"

)

func CreatePurchaseOrderService(db *gorm.DB, poPayload *poModuleModel.PurchaseOrderPayload, roleName string) (string, error) {
	log := logger.InitLogger()
	log.Info("🛠️ CreatePurchaseOrderService invoked")

	// ✅ Step 1: Get current month/year
	now := time.Now()
	month := now.Month()
	year := now.Year() % 100 // last two digits

	// ✅ Step 2: Find the last invoice number for this month/year
	var lastInvoice string
	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
		Select(`"invoiceNumber"`).
		Where(`"invoiceNumber" LIKE ?`, fmt.Sprintf("POINV-%02d%02d-%%", month, year)).
		Order(`purchase_order_id DESC`).
		Limit(1).
		Scan(&lastInvoice).Error

	if err != nil {
		log.Error("❌ Failed to fetch last invoice: " + err.Error())
		return "", err
	}

	// ✅ Step 3: Extract and increment sequence
	sequence := 10001
	if lastInvoice != "" {
		parts := strings.Split(lastInvoice, "-")
		if len(parts) == 3 {
			lastSeq, convErr := strconv.Atoi(parts[2])
			if convErr == nil {
				sequence = lastSeq + 1
			}
		}
	}

	// ✅ Step 4: Build new invoice number
	invoiceNumber := fmt.Sprintf("POINV-%02d%02d-%05d", month, year, sequence)

	// ✅ Step 5: Create Purchase Order
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
		CreatedAt:     now.Format("2006-01-02 15:04:05"),
		CreatedBy:     roleName,
		IsDelete:      false,
		InvoiceNumber: invoiceNumber, // 🧾 Save it here
	}

	if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).Create(&po).Error; err != nil {
		log.Error("❌ Failed to create Purchase Order: " + err.Error())
		return "", err
	}

	log.Infof("✅ Purchase Order created with Invoice: %s (ID: %d)", invoiceNumber, po.PurchaseOrderID)

	// ✅ Step 6: Insert Products
	for _, prod := range poPayload.Products {
		product := poModuleModel.PurchaseOrderProduct{
			PurchaseOrderID: po.PurchaseOrderID,
			CategoryID:      prod.CategoryID,
			Description:     prod.Description,
			UnitPrice:       fmt.Sprintf("%v", prod.UnitPrice),
			Discount:        fmt.Sprintf("%v", prod.Discount),
			Quantity:        fmt.Sprintf("%v", prod.Quantity),
			Total:           fmt.Sprintf("%v", prod.Total),
			CreatedAt:       now.Format("2006-01-02 15:04:05"),
			CreatedBy:       roleName,
		}

		if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts"`).Create(&product).Error; err != nil {
			log.Error("❌ Failed to insert product: " + err.Error())
			return "", err
		}
	}

	// ✅ Transaction logging
	transErr := service.LogTransaction(db, 1, "Admin", 2, fmt.Sprintf("PO Created: %s", invoiceNumber))
	if transErr != nil {
		log.Error("Failed to log transaction : " + transErr.Error())
	} else {
		log.Info("Transaction Log saved Successfully \n\n")
	}

	log.Info("✅ Purchase Order and Products saved successfully")
	return invoiceNumber, nil
}

func GetAllPurchaseOrdersService(db *gorm.DB) []poModuleModel.PurchaseOrderPayload {
	log := logger.InitLogger()

	var purchaseOrders []poModuleModel.PurchaseOrderResponse
	var result []poModuleModel.PurchaseOrderPayload

	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders" AS po`).
		Select(`
        po.purchase_order_id,
        po."invoiceNumber",
        po.sub_total,
        po.total_discount,
        po.tax_enabled,
        po.tax_percentage,
        po.tax_amount,
        po.total_amount,
        po.credited_date,
        po."createdAt",
        po."createdBy",
        po."supplier_id" AS "supplierId",
        s."supplierName" AS "supplierName",
        s."supplierCompanyName" AS "supplierCompany",
        s."supplierCode" AS "supplierCode",
        s."supplierEmail" AS "supplierEmail",
        s."supplierContactNumber" AS "supplierMobile", 
        s."supplierGSTNumber" AS "supplierGST",
        s."supplierPaymentTerms" AS "supplierTerms",
        po."branch_id" AS "branchId",
        b."refBranchName" AS "branchName",
        b."refBranchCode" AS "branchCode",
        b."refLocation" AS "branchLocation",
        b."refMobile" AS "branchMobile",
        b."refEmail" AS "branchEmail",
        b."isMainBranch" AS "isMainBranch",
        b."isActive" AS "isActive"
    `).
		Joins(`LEFT JOIN public."Supplier" s ON po.supplier_id = s."supplierId"`).
		Joins(`LEFT JOIN public."Branches" b ON po.branch_id = b."refBranchId"`).
		Where(`po."isDelete" = ? AND po."invoiceStatus" = ?`, false, false).
		Order("po.purchase_order_id DESC").
		Scan(&purchaseOrders).Error

	if err != nil {
		log.Error("❌ Failed to fetch purchase orders: " + err.Error())
		return result
	}

	for _, po := range purchaseOrders {
		var products []poModuleModel.PurchaseOrderProduct
		db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts"`).
			Where("purchase_order_id = ?", po.PurchaseOrderID).
			Scan(&products)

		for i := range products {
			if products[i].CategoryID != 0 {
				var category poModuleModel.InitialCategory
				err := db.Table(`"public"."InitialCategories"`).
					Where(`"initialCategoryId"= ?`, products[i].CategoryID).
					First(&category).Error
				if err == nil {
					products[i].CategoryDetails = &category
				}
			}
		}

		poPayload := poModuleModel.PurchaseOrderPayload{
			PurchaseOrderID: po.PurchaseOrderID,
			InvoiceNumber:   po.InvoiceNumber,
			Supplier: poModuleModel.SupplierDetails{
				SupplierId:           po.SupplierID,
				SupplierName:         po.SupplierName,
				SupplierCompanyName:  po.SupplierCompany,
				SupplierCode:         po.SupplierCode,
				SupplierEmail:        po.SupplierEmail,
				SupplierMobile:       po.SupplierMobile,
				SupplierGSTNumber:    po.SupplierGST,
				SupplierPaymentTerms: po.SupplierTerms,
			},
			Branch: poModuleModel.BranchDetails{
				RefBranchId:   po.BranchID,
				RefBranchName: po.BranchName,
				RefBranchCode: po.BranchCode,
				RefLocation:   po.BranchLocation,
				RefMobile:     po.BranchMobile,
				RefEmail:      po.BranchEmail,
				IsMainBranch:  po.IsMainBranch,
				IsActive:      po.IsActive,
			},
			Summary: struct {
				SubTotal      string `json:"subTotal"`
				TotalDiscount string `json:"totalDiscount"`
				TaxEnabled    bool   `json:"taxEnabled"`
				TaxPercentage string `json:"taxPercentage"`
				TaxAmount     string `json:"taxAmount"`
				TotalAmount   string `json:"totalAmount"`
			}{
				SubTotal:      po.SubTotal,
				TotalDiscount: po.TotalDiscount,
				TaxEnabled:    po.TaxEnabled,
				TaxPercentage: po.TaxPercentage,
				TaxAmount:     po.TaxAmount,
				TotalAmount:   po.TotalAmount,
			},
			CreditedDate: po.CreditedDate,
			Products:     products,
		}

		result = append(result, poPayload)
	}

	return result
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
