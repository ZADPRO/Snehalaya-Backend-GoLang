package poService

import (
	"errors"
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
		Select(`"purchaseOrderNumber"`).
		Where(`"purchaseOrderNumber" LIKE ?`, fmt.Sprintf("PO-%02d%02d-%%", month, year)).
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
	purchaseOrderNumber := fmt.Sprintf("PO-%02d%02d-%05d", month, year, sequence)

	// ✅ Step 5: Create Purchase Order
	po := poModuleModel.PurchaseOrder{
		SupplierID:          poPayload.Supplier.SupplierId,
		BranchID:            poPayload.Branch.RefBranchId,
		SubTotal:            fmt.Sprintf("%v", poPayload.Summary.SubTotal),
		TotalDiscount:       fmt.Sprintf("%v", poPayload.Summary.TotalDiscount),
		TaxEnabled:          poPayload.Summary.TaxEnabled,
		TaxPercentage:       fmt.Sprintf("%v", poPayload.Summary.TaxPercentage),
		TaxAmount:           fmt.Sprintf("%v", poPayload.Summary.TaxAmount),
		TotalAmount:         fmt.Sprintf("%v", poPayload.Summary.TotalAmount),
		CreditedDate:        poPayload.CreditedDate,
		CreatedAt:           now.Format("2006-01-02 15:04:05"),
		CreatedBy:           roleName,
		IsDelete:            false,
		PurchaseOrderNumber: purchaseOrderNumber, // 🧾 Save it here
	}

	if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).Create(&po).Error; err != nil {
		log.Error("❌ Failed to create Purchase Order: " + err.Error())
		return "", err
	}

	log.Infof("✅ Purchase Order created with Invoice: %s (ID: %d)", purchaseOrderNumber, po.PurchaseOrderID)

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
	transErr := service.LogTransaction(db, 1, "Admin", 2, fmt.Sprintf("PO Created: %s", purchaseOrderNumber))
	if transErr != nil {
		log.Error("Failed to log transaction : " + transErr.Error())
	} else {
		log.Info("Transaction Log saved Successfully \n\n")
	}

	log.Info("✅ Purchase Order and Products saved successfully")
	return purchaseOrderNumber, nil
}

func GetAllPurchaseOrdersService(db *gorm.DB) []poModuleModel.PurchaseOrderPayload {
	log := logger.InitLogger()

	var purchaseOrders []poModuleModel.PurchaseOrderResponse
	var result []poModuleModel.PurchaseOrderPayload

	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders" AS po`).
		Select(`
        po.purchase_order_id,
        po."purchaseOrderNumber",
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
		Where(`po."isDelete" = ?`, false).
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
			PurchaseOrderID:     po.PurchaseOrderID,
			PurchaseOrderNumber: po.PurchaseOrderNumber,
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
func GetAllPurchaseOrdersListService(db *gorm.DB) ([]poModuleModel.PurchaseOrderListResponse, error) {
	log := logger.InitLogger()
	log.Info("📦 GetAllPurchaseOrdersService invoked")

	// Step 1: Fetch all purchase orders (now includes invoice fields)
	orderQuery := `
	SELECT 
		po.purchase_order_id,
		po."purchaseOrderNumber" AS purchase_order_number,
		CASE 
			WHEN po."invoiceStatus" = true THEN 'Approved'
			WHEN po."invoiceStatus" = false THEN 'Created'
			ELSE 'Created'
		END AS status,

		po."invoiceStatus" AS invoice_status,
		po."invoiceFinalNumber" AS invoice_final_number,

		COALESCE(SUM(
			CAST(NULLIF(REGEXP_REPLACE(pop.quantity::text, '[^0-9.-]', '', 'g'), '') AS BIGINT)
		), 0) AS total_ordered_quantity,
		COALESCE(SUM(
			CAST(NULLIF(REGEXP_REPLACE(pop.accepted_quantity::text, '[^0-9.-]', '', 'g'), '') AS BIGINT)
		), 0) AS total_accepted_quantity,
		COALESCE(SUM(
			CAST(NULLIF(REGEXP_REPLACE(pop.rejected_quantity::text, '[^0-9.-]', '', 'g'), '') AS BIGINT)
		), 0) AS total_rejected_quantity,

		po.total_amount,
		po."createdAt" AS created_at, 
		po.tax_amount AS taxable_amount,
		po.supplier_id,
    	s."supplierName" AS supplier_name,
		po.branch_id,
		b."refBranchName" AS branch_name
	FROM "purchaseOrderMgmt"."PurchaseOrders" po
	LEFT JOIN "purchaseOrderMgmt"."PurchaseOrderProducts" pop 
		ON po.purchase_order_id = pop.purchase_order_id
	LEFT JOIN public."Supplier" s 
		ON po.supplier_id = s."supplierId"
	LEFT JOIN public."Branches" b 
		ON po.branch_id = b."refBranchId"
	WHERE (po."isDelete" IS NULL OR po."isDelete" = false)
	GROUP BY 
		po.purchase_order_id, 
		po."purchaseOrderNumber", 
		po."invoiceStatus", 
		po."invoiceFinalNumber",
		po.total_amount, 
		po."createdAt", 
		po.tax_amount, 
		po.supplier_id,
		po.branch_id,
		s."supplierName", 
		b."refBranchName"
	ORDER BY po.purchase_order_id DESC;
	`

	var orders []poModuleModel.PurchaseOrderListResponse
	if err := db.Raw(orderQuery).Scan(&orders).Error; err != nil {
		log.Errorf("❌ Failed to fetch purchase orders: %v", err)
		return nil, err
	}

	// Step 2: For each order, fetch product details
	for i := range orders {
		productQuery := `
SELECT 
    pop.po_product_id,
    pop.purchase_order_id,
    pop.category_id,
    pop.description,
    pop.unit_price,
    pop.discount,
    pop.quantity,
    pop.total,

    -- ✅ Extract numeric value from malformed strings like "%!d(float64=12)"
    COALESCE(
        NULLIF(
            REGEXP_REPLACE(pop.accepted_quantity::text, '.*float64=([0-9.-]+).*', '\1'),
            ''
        ),
        '0'
    ) AS accepted_quantity,

    COALESCE(
        NULLIF(
            REGEXP_REPLACE(pop.rejected_quantity::text, '.*float64=([0-9.-]+).*', '\1'),
            ''
        ),
        '0'
    ) AS rejected_quantity,

    pop."createdAt",
    pop."createdBy",
    pop."updatedAt",
    pop."updatedBy",

    ic."initialCategoryId"   AS initial_category_id,
    ic."initialCategoryName" AS initial_category_name,
    ic."initialCategoryCode" AS initial_category_code,
    ic."isDelete"            AS category_is_delete,
    ic."createdAt"           AS category_created_at,
    ic."createdBy"           AS category_created_by,
    ic."updatedAt"           AS category_updated_at,
    ic."updatedBy"           AS category_updated_by

FROM "purchaseOrderMgmt"."PurchaseOrderProducts" pop
LEFT JOIN public."InitialCategories" ic 
    ON pop.category_id = ic."initialCategoryId"
WHERE pop.purchase_order_id = ?
ORDER BY pop.po_product_id ASC;
`

		var products []poModuleModel.PurchaseOrderProductLatest
		if err := db.Raw(productQuery, orders[i].PurchaseOrderId).Scan(&products).Error; err != nil {
			log.Errorf("❌ Failed to fetch products for PO ID %d: %v", orders[i].PurchaseOrderId, err)
			continue
		}

		// ✅ map flat fields into nested CategoryDetails
		for j := range products {
			products[j].CategoryDetails = poModuleModel.CategoryDetails{
				InitialCategoryId:   products[j].InitialCategoryId,
				InitialCategoryName: products[j].InitialCategoryName,
				InitialCategoryCode: products[j].InitialCategoryCode,
				IsDelete:            products[j].IsDelete,
				CreatedAt:           products[j].CategoryCreatedAt,
				CreatedBy:           products[j].CategoryCreatedBy,
				UpdatedAt:           products[j].CategoryUpdatedAt,
				UpdatedBy:           products[j].CategoryUpdatedBy,
			}
		}

		orders[i].Products = products
	}

	log.Infof("✅ %d Purchase Orders fetched successfully", len(orders))
	return orders, nil
}

type UpdatePOProductRequest struct {
	PurchaseOrderID     int    `json:"purchase_order_id"`
	PurchaseOrderNumber string `json:"purchase_order_number"`
	CategoryID          int    `json:"category_id"`
	POProductID         int    `json:"po_product_id"`
	AcceptedQuantity    string `json:"accepted_quantity"`
	RejectedQuantity    string `json:"rejected_quantity"`
	Status              string `json:"status"`
}

func UpdatePurchaseOrderProductsService(db *gorm.DB, payload []UpdatePOProductRequest) error {
	log := logger.InitLogger()
	log.Info("💾 UpdatePurchaseOrderProductsService invoked")

	tx := db.Begin()

	for _, item := range payload {
		updateQuery := `
			UPDATE "purchaseOrderMgmt"."PurchaseOrderProducts"
			SET 
				accepted_quantity = ?,
				rejected_quantity = ?,
				status = ?,
				"updatedAt" = NOW()
			WHERE po_product_id = ?;
		`

		// ⬇️ Removed fmt.Sprintf conversions here
		if err := tx.Exec(updateQuery,
			item.AcceptedQuantity,
			item.RejectedQuantity,
			item.Status,
			item.POProductID,
		).Error; err != nil {
			tx.Rollback()
			log.Errorf("❌ Failed to update product ID %d: %v", item.POProductID, err)
			return err
		}

		insertInstanceQuery := `
			INSERT INTO "purchaseOrderMgmt"."PurchaseOrderProductInstances"
			(po_product_id, serial_no, category_id, product_description, unit_price, status, "createdAt")
			SELECT 
				pop.po_product_id::varchar,
				CONCAT(pop.po_product_id::varchar, '-', NOW()),
				pop.category_id,
				pop.description,
				pop.unit_price::varchar,
				?,
				TO_CHAR(NOW(), 'YYYY-MM-DD HH24:MI:SS')
			FROM "purchaseOrderMgmt"."PurchaseOrderProducts" pop
			WHERE pop.po_product_id = ?;
		`

		if err := tx.Exec(insertInstanceQuery, item.Status, strconv.Itoa(item.POProductID)).Error; err != nil {
			tx.Rollback()
			log.Errorf("❌ Failed to insert instance for product ID %d: %v", item.POProductID, err)
			return err
		}

		log.Infof("✅ Product %d updated successfully", item.POProductID)
	}

	if err := tx.Commit().Error; err != nil {
		log.Errorf("❌ Commit failed: %v", err)
		return err
	}

	log.Info("✅ All products updated successfully")
	return nil
}

type SavePurchaseOrderProductsRequest struct {
	PurchaseOrderId int                    `json:"purchaseOrderId"`
	Products        []SavePOProductRequest `json:"products"`
}

type SavePOProductRequest struct {
	BranchId      int                      `json:"productBranchId"`
	SNo           int                      `json:"sNo"`
	LineNumber    int                      `json:"lineNumber"`
	ProductName   string                   `json:"productName"`
	Brand         string                   `json:"brand"`
	CategoryId    int                      `json:"categoryId"`
	SubCategoryId int                      `json:"subCategoryId"`
	TaxClass      string                   `json:"taxClass"`
	Quantity      int                      `json:"quantity"`
	Cost          float64                  `json:"cost"`
	ProfitMargin  float64                  `json:"profitMargin"`
	SellingPrice  float64                  `json:"sellingPrice"`
	MRP           float64                  `json:"mrp"`
	DialogRows    []SavePODialogRowRequest `json:"dialogRows"`
}

type SavePODialogRowRequest struct {
	BranchId           int     `json:"productBranchId"`
	SNo                int     `json:"sNo"`
	POProductID        int     `json:"poProductId"`
	LineNumber         int     `json:"lineNumber"`
	ReferenceNumber    string  `json:"referenceNumber"`
	ProductDescription string  `json:"productDescription"`
	Discount           float64 `json:"discount"`
	Price              float64 `json:"price"`
	DiscountPrice      float64 `json:"discountPrice"`
	Margin             float64 `json:"margin"`
	TotalAmount        string  `json:"totalAmount"`
}

func SavePurchaseOrderProductsService(db *gorm.DB, payload SavePurchaseOrderProductsRequest) error {
	type DialogRow struct {
		BranchId           int     `json:"productBranchId"`
		SNo                int     `json:"sNo"`
		LineNumber         int     `json:"lineNumber"`
		ReferenceNumber    string  `json:"referenceNumber"`
		ProductDescription string  `json:"productDescription"`
		Discount           float64 `json:"discount"`
		Price              float64 `json:"price"`
		DiscountPrice      float64 `json:"discountPrice"`
		Margin             float64 `json:"margin"`
		TotalAmount        string  `json:"totalAmount"`
	}

	type Product struct {
		BranchId      int         `json:"productBranchId"`
		SNo           int         `json:"sNo"`
		LineNumber    int         `json:"lineNumber"`
		ProductName   string      `json:"productName"`
		Brand         string      `json:"brand"`
		CategoryId    int         `json:"categoryId"`
		SubCategoryId int         `json:"subCategoryId"`
		TaxClass      string      `json:"taxClass"`
		Quantity      int         `json:"quantity"`
		Cost          float64     `json:"cost"`
		ProfitMargin  float64     `json:"profitMargin"`
		SellingPrice  float64     `json:"sellingPrice"`
		Mrp           float64     `json:"mrp"`
		DialogRows    []DialogRow `json:"dialogRows"`
	}

	type PurchaseOrderAcceptedProduct struct {
		ProductBranchId    int    `gorm:"column:productBranchId"`
		ProductInstanceId  int    `gorm:"primaryKey;autoIncrement;column:product_instance_id"`
		PoProductId        int    `gorm:"column:po_product_id"`
		PurchaseOrderId    int    `gorm:"column:purchaseOrderId"`
		LineNumber         string `gorm:"column:line_number"`
		ReferenceNumber    string `gorm:"column:reference_number"`
		ProductDescription string `gorm:"column:product_description"`
		Discount           string `gorm:"column:discount"`
		UnitPrice          string `gorm:"column:unit_price"`
		DiscountPrice      string `gorm:"column:discount_price"`
		Margin             string `gorm:"column:margin"`
		TotalAmount        string `gorm:"column:total_amount"`
		CategoryId         int    `gorm:"column:category_id"`
		SubCategoryId      int    `gorm:"column:sub_category_id"`
		Status             string `gorm:"column:status"`
		CreatedAt          string `gorm:"column:createdAt"`
		CreatedBy          string `gorm:"column:createdBy"`
		UpdatedAt          string `gorm:"column:updatedAt"`
		UpdatedBy          string `gorm:"column:updatedBy"`
		IsDelete           bool   `gorm:"column:isDelete"`
		ProductName        string `gorm:"column:product_name"`
		SKU                string `gorm:"column:SKU"`
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	monthYear := time.Now().Format("0106") // e.g. 1025 for Oct 2025

	// --- Determine latest SKU number ---
	var lastSKU string
	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts"`).
		Select(`"SKU"`).
		Order("product_instance_id DESC").
		Limit(1).
		Scan(&lastSKU).Error

	var nextSKUCode int64 = 10001
	if err == nil && lastSKU != "" {
		parts := strings.Split(lastSKU, "-")
		if len(parts) >= 3 {
			if num, convErr := strconv.ParseInt(parts[2], 10, 64); convErr == nil {
				nextSKUCode = num + 1
			}
		}
	}

	var records []PurchaseOrderAcceptedProduct
	for _, product := range payload.Products {
		for _, row := range product.DialogRows {
			sku := fmt.Sprintf("SS-%s-%05d", monthYear, nextSKUCode)
			nextSKUCode++ // increment for next product

			record := PurchaseOrderAcceptedProduct{
				ProductBranchId:    product.BranchId,
				PoProductId:        payload.PurchaseOrderId,
				PurchaseOrderId:    payload.PurchaseOrderId,
				LineNumber:         fmt.Sprintf("%d", row.LineNumber),
				ReferenceNumber:    row.ReferenceNumber,
				ProductDescription: row.ProductDescription,
				Discount:           fmt.Sprintf("%v", row.Discount),
				UnitPrice:          fmt.Sprintf("%v", row.Price),
				DiscountPrice:      fmt.Sprintf("%v", row.DiscountPrice),
				Margin:             fmt.Sprintf("%v", row.Margin),
				TotalAmount:        row.TotalAmount,
				CategoryId:         product.CategoryId,
				SubCategoryId:      product.SubCategoryId,
				ProductName:        product.ProductName,
				Status:             "Active",
				CreatedAt:          currentTime,
				CreatedBy:          "Admin",
				UpdatedAt:          currentTime,
				UpdatedBy:          "Admin",
				IsDelete:           false,
				SKU:                sku,
			}

			records = append(records, record)
		}
	}

	if len(records) == 0 {
		return fmt.Errorf("no dialog rows to insert")
	}

	// --- Bulk insert accepted products ---
	if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts"`).Create(&records).Error; err != nil {
		return fmt.Errorf("failed to insert accepted products: %v", err)
	}

	// --- Generate Invoice Number ---
	var lastInvoice string
	err = db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
		Select(`"invoiceFinalNumber"`).
		Order("purchase_order_id DESC").
		Limit(1).
		Scan(&lastInvoice).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to fetch last invoice: %v", err)
	}

	var nextNumber int64 = 10001
	if lastInvoice != "" {
		parts := strings.Split(lastInvoice, "-")
		if len(parts) >= 4 {
			if num, convErr := strconv.ParseInt(parts[3], 10, 64); convErr == nil {
				nextNumber = num + 1
			}
		}
	}

	invoiceNumber := fmt.Sprintf("PO-INV-%s-%05d", monthYear, nextNumber)

	updateData := map[string]interface{}{
		`"invoiceStatus"`:      true,
		`"invoiceFinalNumber"`: invoiceNumber,
		`"updatedAt"`:          currentTime,
		`"updatedBy"`:          "Admin",
	}

	if err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
		Where("purchase_order_id = ?", payload.PurchaseOrderId).
		Updates(updateData).Error; err != nil {
		return fmt.Errorf("failed to update PurchaseOrders with invoice info: %v", err)
	}

	return nil
}

// --- Structs ---
type DialogRow struct {
	SNo                int     `json:"sNo"`
	LineNumber         int     `json:"lineNumber"`
	ReferenceNumber    string  `json:"referenceNumber"`
	ProductDescription string  `json:"productDescription"`
	Discount           float64 `json:"discount"`
	Price              float64 `json:"price"`
	DiscountPrice      float64 `json:"discountPrice"`
	Margin             float64 `json:"margin"`
	TotalAmount        string  `json:"totalAmount"`
}

type Product struct {
	SNo             int         `json:"sNo"`
	LineNumber      int         `json:"lineNumber"`
	ProductName     string      `json:"productName"`
	Brand           string      `json:"brand"`
	CategoryId      int         `json:"categoryId"`
	SubCategoryId   int         `json:"subCategoryId"`
	TaxClass        string      `json:"taxClass"`
	Quantity        int         `json:"quantity"`
	Cost            float64     `json:"cost"`
	ProfitMargin    float64     `json:"profitMargin"`
	SellingPrice    float64     `json:"sellingPrice"`
	Mrp             float64     `json:"mrp"`
	DiscountPercent float64     `json:"discountPercent"`
	DiscountPrice   float64     `json:"discountPrice"`
	DialogRows      []DialogRow `json:"dialogRows"`
}

type PurchaseOrderDetailsResponse struct {
	PurchaseOrderId    int       `json:"purchaseOrderId"`
	InvoiceFinalNumber string    `json:"invoiceFinalNumber"`
	Products           []Product `json:"products"`
}

func GetPurchaseOrderDetailsService(db *gorm.DB, purchaseOrderNumber string) (PurchaseOrderDetailsResponse, error) {
	log := logger.InitLogger()
	var response PurchaseOrderDetailsResponse

	// Step 1: Get purchase order basic info (ID + invoice number)
	var po struct {
		PurchaseOrderId    int    `gorm:"column:purchase_order_id"`
		InvoiceFinalNumber string `gorm:"column:invoiceFinalNumber"`
	}
	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
		Select(`"purchase_order_id", "invoiceFinalNumber"`).
		Where(`"purchaseOrderNumber" = ?`, purchaseOrderNumber).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response, fmt.Errorf("purchase order not found")
	}
	if err != nil {
		return response, err
	}

	// Step 2: Fetch PO products (join with categories/subcategories)
	type ProductRow struct {
		PoProductId     int     `gorm:"column:po_product_id"`
		LineNumber      int     `gorm:"column:line_number"`
		Description     string  `gorm:"column:description"`
		CategoryId      int     `gorm:"column:category_id"`
		CategoryName    string  `gorm:"column:category_name"`
		SubCategoryId   int     `gorm:"column:sub_category_id"`
		SubCategoryName string  `gorm:"column:sub_category_name"`
		Quantity        int     `gorm:"column:quantity"`
		UnitPrice       float64 `gorm:"column:unit_price"`
	}

	var poProducts []ProductRow
	err = db.Table(`"purchaseOrderMgmt"."PurchaseOrderProducts" AS pop`).
		Select(`
			pop."po_product_id",
			pop."line_number",
			pop."description",
			pop."category_id",
			c."categoryName" AS category_name,
			pop."sub_category_id",
			sc."subCategoryName" AS sub_category_name,
			pop."quantity",
			pop."unit_price"
		`).
		Joins(`LEFT JOIN public."Categories" c ON c."refCategoryid" = pop."category_id"`).
		Joins(`LEFT JOIN public."SubCategories" sc ON sc."refSubCategoryId" = pop."sub_category_id"`).
		Where(`pop."purchase_order_id" = ?`, po.PurchaseOrderId).
		Scan(&poProducts).Error
	if err != nil {
		return response, err
	}

	// Step 3: Fetch accepted products (dialog rows + category/subcategory)
	type AcceptedRow struct {
		PoProductId        int     `gorm:"column:po_product_id"`
		LineNumber         int     `gorm:"column:line_number"`
		ReferenceNumber    string  `gorm:"column:reference_number"`
		ProductDescription string  `gorm:"column:product_description"`
		Discount           float64 `gorm:"column:discount"`
		UnitPrice          float64 `gorm:"column:unit_price"`
		DiscountPrice      float64 `gorm:"column:discount_price"`
		Margin             float64 `gorm:"column:margin"`
		TotalAmount        string  `gorm:"column:total_amount"`
		CategoryId         int     `gorm:"column:category_id"`
		CategoryName       string  `gorm:"column:category_name"`
		SubCategoryId      int     `gorm:"column:sub_category_id"`
		SubCategoryName    string  `gorm:"column:sub_category_name"`
		ProductName        string  `gorm:"column:product_name"`
		SKU                string  `gorm:"column:SKU"`
	}

	var acceptedProducts []AcceptedRow
	err = db.Table(`"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts" AS ap`).
		Select(`
			ap."po_product_id",
			ap."line_number",
			ap."reference_number",
			ap."product_description",
			ap."discount",
			ap."unit_price",
			ap."discount_price",
			ap."margin",
			ap."total_amount",
			ap."category_id",
			c."categoryName" AS category_name,
			ap."sub_category_id",
			sc."subCategoryName" AS sub_category_name,
			ap."product_name",
			ap."SKU"
		`).
		Joins(`LEFT JOIN public."Categories" c ON c."refCategoryid" = ap."category_id"`).
		Joins(`LEFT JOIN public."SubCategories" sc ON sc."refSubCategoryId" = ap."sub_category_id"`).
		Where(`ap."purchaseOrderId" = ?`, po.PurchaseOrderId).
		Order(`ap."product_instance_id" ASC`).
		Scan(&acceptedProducts).Error
	if err != nil {
		return response, err
	}

	// Step 4: Combine both into final response
	var responseProducts []Product
	sNo := 1
	for _, p := range poProducts {
		var rows []DialogRow
		rowIndex := 1
		for _, ap := range acceptedProducts {
			if ap.PoProductId == p.PoProductId {
				rows = append(rows, DialogRow{
					SNo:                rowIndex,
					LineNumber:         ap.LineNumber,
					ReferenceNumber:    ap.ReferenceNumber,
					ProductDescription: ap.ProductDescription,
					Discount:           ap.Discount,
					Price:              ap.UnitPrice,
					DiscountPrice:      ap.DiscountPrice,
					Margin:             ap.Margin,
					TotalAmount:        ap.TotalAmount,
				})
				rowIndex++
			}
		}

		responseProducts = append(responseProducts, Product{
			SNo:             sNo,
			LineNumber:      p.LineNumber,
			ProductName:     p.Description,
			Brand:           "Snehalayaa",
			CategoryId:      p.CategoryId,
			SubCategoryId:   p.SubCategoryId,
			TaxClass:        "HSN Code",
			Quantity:        p.Quantity,
			Cost:            p.UnitPrice,
			ProfitMargin:    80.5,
			SellingPrice:    p.UnitPrice * 1.8,
			Mrp:             p.UnitPrice * 1.8,
			DiscountPercent: 0,
			DiscountPrice:   0,
			DialogRows:      rows,
		})
		sNo++
	}

	response = PurchaseOrderDetailsResponse{
		PurchaseOrderId:    po.PurchaseOrderId,
		InvoiceFinalNumber: po.InvoiceFinalNumber,
		Products:           responseProducts,
	}

	log.Info(fmt.Sprintf("✅ Loaded PO #%s (%s) with %d products", purchaseOrderNumber, po.InvoiceFinalNumber, len(responseProducts)))
	return response, nil
}

func GetAcceptedProductsService(db *gorm.DB, purchaseOrderId string) ([]map[string]interface{}, error) {
	log := logger.InitLogger()
	log.Infof("🔍 Fetching accepted products for PurchaseOrderId: %s", purchaseOrderId)

	// --- Define model mapping DB columns ---
	type PurchaseOrderAcceptedProduct struct {
		ProductInstanceId  int     `json:"product_instance_id" gorm:"column:product_instance_id"`
		PoProductId        int     `json:"po_product_id" gorm:"column:po_product_id"`
		PurchaseOrderId    int     `json:"purchaseOrderId" gorm:"column:purchaseOrderId"`
		LineNumber         string  `json:"line_number" gorm:"column:line_number"`
		ReferenceNumber    string  `json:"reference_number" gorm:"column:reference_number"`
		ProductDescription string  `json:"product_description" gorm:"column:product_description"`
		Discount           string  `json:"discount" gorm:"column:discount"`
		UnitPrice          string  `json:"unit_price" gorm:"column:unit_price"`
		DiscountPrice      string  `json:"discount_price" gorm:"column:discount_price"`
		Margin             string  `json:"margin" gorm:"column:margin"`
		TotalAmount        string  `json:"total_amount" gorm:"column:total_amount"`
		CategoryId         int     `json:"category_id" gorm:"column:category_id"`
		SubCategoryId      int     `json:"sub_category_id" gorm:"column:sub_category_id"`
		ProductName        string  `json:"product_name" gorm:"column:product_name"`
		SKU                string  `json:"SKU" gorm:"column:SKU"`
		Status             string  `json:"status" gorm:"column:status"`
		CreatedAt          string  `json:"createdAt" gorm:"column:createdAt"`
		CreatedBy          string  `json:"createdBy" gorm:"column:createdBy"`
		UpdatedAt          string  `json:"updatedAt" gorm:"column:updatedAt"`
		UpdatedBy          string  `json:"updatedBy" gorm:"column:updatedBy"`
		IsDelete           bool    `json:"isDelete" gorm:"column:isDelete"`
		ProductBranchId    *int    `json:"productBranchId" gorm:"column:productBranchId"`
		Quantity           *string `json:"quantity" gorm:"column:quantity"`
	}

	var records []PurchaseOrderAcceptedProduct

	// --- Fetch only active (non-deleted) records ---
	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts"`).
		Where(`"purchaseOrderId" = ? AND "isDelete" = false`, purchaseOrderId).
		Order(`"product_instance_id" ASC`).
		Find(&records).Error

	if err != nil {
		log.Errorf("❌ Database query failed for PurchaseOrderId %s: %v", purchaseOrderId, err)
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	if len(records) == 0 {
		log.Warnf("⚠️ No accepted products found for PurchaseOrderId: %s", purchaseOrderId)
		return []map[string]interface{}{}, nil
	}

	// --- Prepare final JSON-friendly result ---
	var result []map[string]interface{}
	for _, rec := range records {
		// Only include non-deleted records (extra safeguard)
		if !rec.IsDelete {
			result = append(result, map[string]interface{}{
				"productInstanceId":  rec.ProductInstanceId,
				"poProductId":        rec.PoProductId,
				"purchaseOrderId":    rec.PurchaseOrderId,
				"lineNumber":         rec.LineNumber,
				"referenceNumber":    rec.ReferenceNumber,
				"productDescription": rec.ProductDescription,
				"discount":           rec.Discount,
				"unitPrice":          rec.UnitPrice,
				"discountPrice":      rec.DiscountPrice,
				"margin":             rec.Margin,
				"totalAmount":        rec.TotalAmount,
				"categoryId":         rec.CategoryId,
				"subCategoryId":      rec.SubCategoryId,
				"productName":        rec.ProductName,
				"SKU":                rec.SKU,
				"status":             rec.Status,
				"createdAt":          rec.CreatedAt,
				"createdBy":          rec.CreatedBy,
				"updatedAt":          rec.UpdatedAt,
				"updatedBy":          rec.UpdatedBy,
				"productBranchId":    rec.ProductBranchId,
				"quantity":           rec.Quantity,
			})
		}
	}

	log.Infof("📦 Successfully fetched %d accepted products for PurchaseOrderId: %s", len(result), purchaseOrderId)
	return result, nil
}

type PurchaseOrderAcceptedProductResponse struct {
	ProductInstanceID  int    `json:"productInstanceId" gorm:"column:product_instance_id"`
	PoProductID        int    `json:"poProductId" gorm:"column:po_product_id"`
	LineNumber         string `json:"lineNumber" gorm:"column:line_number"`
	ReferenceNumber    string `json:"referenceNumber" gorm:"column:reference_number"`
	ProductDescription string `json:"productDescription" gorm:"column:product_description"`
	Discount           string `json:"discount" gorm:"column:discount"`
	UnitPrice          string `json:"unitPrice" gorm:"column:unit_price"`
	DiscountPrice      string `json:"discountPrice" gorm:"column:discount_price"`
	Margin             string `json:"margin" gorm:"column:margin"`
	TotalAmount        string `json:"totalAmount" gorm:"column:total_amount"`
	CategoryID         int    `json:"categoryId" gorm:"column:category_id"`
	SubCategoryID      int    `json:"subCategoryId" gorm:"column:sub_category_id"`
	Status             string `json:"status" gorm:"column:status"`
	CreatedAt          string `json:"createdAt" gorm:"column:created_at"`
	CreatedBy          string `json:"createdBy" gorm:"column:created_by"`
	UpdatedAt          string `json:"updatedAt" gorm:"column:updated_at"`
	UpdatedBy          string `json:"updatedBy" gorm:"column:updated_by"`
	IsDelete           bool   `json:"isDelete" gorm:"column:is_delete"`
	ProductName        string `json:"productName" gorm:"column:product_name"`
	PurchaseOrderId    int    `json:"purchaseOrderId" gorm:"column:purchase_order_id"`
	SKU                string `json:"sku" gorm:"column:sku"`
	ProductBranchId    int    `json:"productBranchId" gorm:"column:product_branch_id"`
	Quantity           string `json:"quantity" gorm:"column:quantity"`
	InvoiceFinalNumber string `json:"invoiceFinalNumber" gorm:"column:invoice_final_number"`
	CategoryName       string `json:"categoryName" gorm:"column:category_name"`
	SubCategoryName    string `json:"subCategoryName" gorm:"column:sub_category_name"`
	BranchName         string `json:"branchName" gorm:"column:branch_name"`
}

func GetAllPurchaseOrderAcceptedProductsService(db *gorm.DB) []PurchaseOrderAcceptedProductResponse {
	log := logger.InitLogger()
	log.Info("💾 GetAllPurchaseOrderAcceptedProductsService invoked")

	var results []PurchaseOrderAcceptedProductResponse

	rawQuery := `
		SELECT
		ap.product_instance_id AS product_instance_id,
		ap.po_product_id AS po_product_id,
		ap.line_number AS line_number,
		ap.reference_number AS reference_number,
		ap.product_description AS product_description,
		ap.discount AS discount,
		ap.unit_price AS unit_price,
		ap.discount_price AS discount_price,
		ap.margin AS margin,
		ap.total_amount AS total_amount,
		ap.category_id AS category_id,
		ap.sub_category_id AS sub_category_id,
		ap.status AS status,
		ap."createdAt" AS created_at,
		ap."createdBy" AS created_by,
		ap."updatedAt" AS updated_at,
		ap."updatedBy" AS updated_by,
		ap."isDelete" AS is_delete,
		ap.product_name AS product_name,
		ap."purchaseOrderId" AS purchase_order_id,
		ap."SKU" AS sku,
		ap."productBranchId" AS product_branch_id,
		ap.quantity AS quantity,
		po."invoiceFinalNumber" AS invoice_final_number,
		c."categoryName" AS category_name,
		sc."subCategoryName" AS sub_category_name,
		b."refBranchName" AS branch_name
		FROM
		"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts" AS ap
		LEFT JOIN "purchaseOrderMgmt"."PurchaseOrders" po ON ap."purchaseOrderId" = po.purchase_order_id
		LEFT JOIN public."Categories" c ON c."refCategoryid" = ap.category_id
		LEFT JOIN public."SubCategories" sc ON sc."refSubCategoryId" = ap.sub_category_id
		LEFT JOIN public."Branches" b ON b."refBranchId" = ap."productBranchId"
		WHERE
		ap."isDelete" = false
		ORDER BY
		ap.product_instance_id DESC;
	`

	err := db.Raw(rawQuery).Scan(&results).Error
	fmt.Print("\n\n\n\nresults", results)
	if err != nil {
		log.Errorf("❌ Failed to fetch accepted products: %v", err)
		return nil
	}

	log.Infof("✅ Retrieved %d accepted products", len(results))
	return results
}

func GetPurchaseOrderFullDetailsService(db *gorm.DB, purchaseOrderNumber string) (PurchaseOrderDetailsResponse, error) {
	log := logger.InitLogger()
	var response PurchaseOrderDetailsResponse

	// Step 1: Fetch PO ID and Invoice using raw SQL
	var po struct {
		PurchaseOrderId    int
		InvoiceFinalNumber string
	}
	poQuery := `
SELECT
    po."purchase_order_id",
    po."invoiceFinalNumber"
FROM
    "purchaseOrderMgmt"."PurchaseOrders" po
WHERE
    po."purchase_order_id" = ?
ORDER BY
    po."purchase_order_id" ASC
LIMIT 1;
`

	err := db.Raw(poQuery, purchaseOrderNumber).Scan(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || po.PurchaseOrderId == 0 {
		return response, fmt.Errorf("purchase order not found")
	}
	if err != nil {
		return response, err
	}

	response.PurchaseOrderId = po.PurchaseOrderId
	response.InvoiceFinalNumber = po.InvoiceFinalNumber

	// Step 2: Fetch all accepted products (dialog rows)
	var acceptedProducts []Product
	productQuery := `
SELECT 
    ap.po_product_id,
    ap.line_number,
    ap.reference_number,
    ap.product_description,
    ap.discount,
    ap.unit_price,
    ap.discount_price,
    ap.margin,
    ap.total_amount,
    ap.category_id,
    c."categoryName" AS category_name,
    ap.sub_category_id,
    sc."subCategoryName" AS sub_category_name,
    ap.product_name,
    ap."SKU"
FROM "purchaseOrderMgmt"."PurchaseOrderAcceptedProducts" AS ap
LEFT JOIN public."Categories" AS c ON c."refCategoryid" = ap.category_id
LEFT JOIN public."SubCategories" AS sc ON sc."refSubCategoryId" = ap.sub_category_id
WHERE ap.po_product_id = $1
ORDER BY ap.po_product_id ASC;
`

	var rows []struct {
		PoProductId        int
		LineNumber         int
		ReferenceNumber    string
		ProductDescription string
		Discount           float64
		UnitPrice          float64
		DiscountPrice      float64
		Margin             float64
		TotalAmount        string
		CategoryId         int
		CategoryName       string
		SubCategoryId      int
		SubCategoryName    string
		ProductName        string
		SKU                string
	}

	if err := db.Raw(productQuery, po.PurchaseOrderId).Scan(&rows).Error; err != nil {
		log.Error("❌ Failed to fetch accepted products: " + err.Error())
		return response, err
	}

	// Map rows into Product + DialogRows
	productMap := map[int]*Product{}
	for _, row := range rows {
		if _, exists := productMap[row.PoProductId]; !exists {
			productMap[row.PoProductId] = &Product{
				ProductName:   row.ProductName,
				CategoryId:    row.CategoryId,
				SubCategoryId: row.SubCategoryId,
				DialogRows:    []DialogRow{},
			}
		}
		productMap[row.PoProductId].DialogRows = append(productMap[row.PoProductId].DialogRows, DialogRow{
			LineNumber:         row.LineNumber,
			ReferenceNumber:    row.ReferenceNumber,
			ProductDescription: row.ProductDescription,
			Discount:           row.Discount,
			Price:              row.UnitPrice,
			DiscountPrice:      row.DiscountPrice,
			Margin:             row.Margin,
			TotalAmount:        row.TotalAmount,
		})
	}

	for _, p := range productMap {
		acceptedProducts = append(acceptedProducts, *p)
	}

	response.Products = acceptedProducts
	return response, nil
}
