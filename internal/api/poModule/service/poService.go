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

	// Step 1: Fetch all purchase orders (same as before)
	orderQuery := `
	SELECT 
		po.purchase_order_id,
		po."purchaseOrderNumber" AS purchase_order_number,
		CASE 
			WHEN po."invoiceStatus" = true THEN 'Approved'
			WHEN po."invoiceStatus" = false THEN 'Created'
			ELSE 'Created'
		END AS status,
		COALESCE(SUM(CAST(pop.quantity AS BIGINT)), 0) AS total_ordered_quantity,
		COALESCE(SUM(CAST(pop.accepted_quantity AS BIGINT)), 0) AS total_accepted_quantity,
		COALESCE(SUM(CAST(pop.rejected_quantity AS BIGINT)), 0) AS total_rejected_quantity,
		po.total_amount,
		po."createdAt" AS created_at, 
		po.tax_amount AS taxable_amount,
		po.supplier_id,
    	s."supplierName" AS supplier_name,                    -- ✅ snake_case alias
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
	PurchaseOrderID     int     `json:"purchase_order_id"`
	PurchaseOrderNumber string  `json:"purchase_order_number"`
	CategoryID          int     `json:"category_id"`
	POProductID         int     `json:"po_product_id"`
	AcceptedQuantity    float64 `json:"accepted_quantity"`
	RejectedQuantity    float64 `json:"rejected_quantity"`
	Status              string  `json:"status"`
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

		if err := tx.Exec(updateQuery, fmt.Sprintf("%d", item.AcceptedQuantity),
			fmt.Sprintf("%d", item.RejectedQuantity), item.Status, item.POProductID).Error; err != nil {
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
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	var records []PurchaseOrderAcceptedProduct

	// --- Build bulk insert records ---
	for _, product := range payload.Products {
		for _, row := range product.DialogRows {
			record := PurchaseOrderAcceptedProduct{
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
	// Get the latest invoice number suffix from PurchaseOrders table
	var lastInvoice string
	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrders"`).
		Select(`"invoiceFinalNumber"`).
		Order("purchase_order_id DESC").
		Limit(1).
		Scan(&lastInvoice).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to fetch last invoice: %v", err)
	}

	// Determine the next invoice number suffix
	var nextNumber int64 = 10001
	if lastInvoice != "" {
		parts := strings.Split(lastInvoice, "-")
		if len(parts) >= 4 {
			if num, convErr := strconv.ParseInt(parts[3], 10, 64); convErr == nil {
				nextNumber = num + 1
			}
		}
	}

	monthYear := time.Now().Format("0106") // e.g., 1025 for Oct 2025
	invoiceNumber := fmt.Sprintf("PO-INV-%s-%05d", monthYear, nextNumber)

	// --- Update PurchaseOrder table ---
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
