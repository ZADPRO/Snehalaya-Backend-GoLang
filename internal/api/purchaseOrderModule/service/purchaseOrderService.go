package purchaseOrderService

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	transactionLogger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/service"
	purchaseOrderModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/model"
	purchaseOrderQuery "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/query"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreatePurchaseOrderService(db *gorm.DB, payload *purchaseOrderModel.CreatePORequest, createdBy string) error {

	poNumber, err := purchaseOrderQuery.GeneratePONumber(db)
	if err != nil {
		return fmt.Errorf("failed to generate PO number: %v", err)
	}
	// 1. Insert into CreatePurchaseOrder
	order := purchaseOrderModel.CreatePurchaseOrder{
		PONumber:        poNumber,
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
		IsInternalPO:    payload.IsInternalPO,
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
		Order(`"dummyProductsId" ASC`).
		Find(&dummyProducts).Error; err != nil {
		return nil, err
	}

	fmt.Printf("\n\n\n\n\ndummyProducts => %v\n\n", dummyProducts)

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

	// Save all products
	for _, product := range products {
		if err := db.Save(&product).Error; err != nil {
			return err
		}
	}
	return nil
}

type ReceivedDummyProductWithPO struct {
	purchaseOrderModel.ProductsDummyAcceptance
	PONumber string `gorm:"column:poNumber" json:"poNumber"` // ensure column tag matches SQL
}

type ReceivedDummyProduct struct {
	Name    string `json:"name" gorm:"column:name"`
	HSNCode string `json:"hsnCode" gorm:"column:HSNCode"`
	SKU     string `json:"sku" gorm:"column:sku"`
	Price   string `json:"price" gorm:"column:price"`
	Status  string `json:"status" gorm:"column:status"`
}

func GetReceivedDummyProductsService(db *gorm.DB) ([]ReceivedDummyProductWithPO, error) {
	var result []ReceivedDummyProductWithPO

	query := `
		SELECT 
			"PDA".*, 
			"CPO"."poNumber"
		FROM 
			"purchaseOrder"."ProductsDummyAcceptance" AS "PDA"
		JOIN 
			"purchaseOrder"."CreatePurchaseOrder" AS "CPO"
		ON 
			"PDA"."purchaseOrderId" = "CPO"."purchaseOrderId"
		WHERE 
			"PDA"."isReceived" = 'true' AND 
			"PDA"."isDelete" = 'false' AND 
			"PDA"."acceptanceStatus" = 'Received'
		ORDER BY 
			"PDA"."dummyProductsId" ASC;
	`

	err := db.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetReceivedDummyProductsBarcodeService(db *gorm.DB) ([]ReceivedDummyProduct, error) {
	var result []ReceivedDummyProduct

	query := `
		SELECT 
			p."name",
			pda."HSNCode",
			pda."dummySKU" AS sku,
			pda."price",
			'Created' AS status
		FROM 
			"purchaseOrder"."ProductsDummyAcceptance" pda
		JOIN 
			"purchaseOrder".products p 
				ON pda."dummySKU" = p.sku
		WHERE 
			(pda."isDelete" IS NULL OR pda."isDelete" != 'true')
			AND pda."acceptanceStatus" = 'Created'
			AND pda."isReceived" = 'true';
	`

	err := db.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	fmt.Println("Result:", result)
	return result, nil
}

func CreateProductService(db *gorm.DB, product *purchaseOrderModel.Product) error {
	log := logger.InitLogger()

	// Check for duplicate SKU
	var existing purchaseOrderModel.Product
	err := db.Table(`"purchaseOrder".products`).
		Where(`sku = ? AND "isDelete" = ?`, product.SKU, "false").
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate product SKU found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("DB error while checking for duplicates: " + err.Error())
		return err
	}

	// Set audit fields
	now := time.Now().Format("2006-01-02 15:04:05")
	product.CreatedAt = now
	product.UpdatedAt = now
	product.CreatedBy = "Admin"
	product.UpdatedBy = "Admin"
	product.IsDelete = "false"

	// Insert into DB
	err = db.Table(`"purchaseOrder".products`).Create(product).Error
	if err != nil {
		log.Error("Failed to create product: " + err.Error())
		return err
	}

	// ✅ Update the acceptanceStatus to "Created"
	err = db.Table(`"purchaseOrder"."ProductsDummyAcceptance"`).
		Where(`"dummySKU" = ?`, product.SKU).
		Update(`acceptanceStatus`, "Created").Error
	if err != nil {
		log.Error("Failed to update acceptanceStatus: " + err.Error())
		return err
	}

	return nil
}

type PurchaseOrderItem struct {
	CategoryId         int     `json:"categoryId"`
	SubCategoryId      int     `json:"subCategoryId"`
	ProductDescription string  `json:"productDescription"`
	UnitPrice          float64 `json:"unitPrice"`
	Quantity           float64 `json:"quantity"`
	DiscountPercent    float64 `json:"discountPercent"`
	DiscountAmount     float64 `json:"discountAmount"`
	Total              float64 `json:"total"`
}

type PurchaseOrderPayload struct {
	SupplierId  int                 `json:"supplierId"`
	BranchId    int                 `json:"branchId"`
	TaxEnabled  bool                `json:"taxEnabled"`
	TaxRate     float64             `json:"taxRate"`
	PaymentFee  float64             `json:"paymentFee"`
	ShippingFee float64             `json:"shippingFee"`
	Subtotal    float64             `json:"subtotal"`
	TaxAmount   float64             `json:"taxAmount"`
	RoundOff    float64             `json:"roundOff"`
	Total       float64             `json:"total"`
	Items       []PurchaseOrderItem `json:"items"`
}

func GeneratePONumber(db *gorm.DB, year int, month int) (string, error) {
	log := logger.InitLogger()
	log.Info("🧮 Generating new PO Number...")

	// A = 2025, B = 2026, C = 2027 ...
	yearCode := string(rune('A' + (year - 2025)))

	// A = Jan, B = Feb ... K = Nov, L = Dec
	monthCode := string(rune('A' + (month - 1)))

	prefix := fmt.Sprintf("PO%s%s", yearCode, monthCode)

	log.Infof("🔍 PO Prefix = %s", prefix)

	// ==== FETCH LATEST PO NUMBER WITH THIS PREFIX ====
	var lastPONumber string

	err := db.
		Raw(`
            SELECT po_number 
            FROM "PurchaseOrderManagement"."PurchaseOrders"
            WHERE po_number LIKE ?
            ORDER BY id DESC 
            LIMIT 1
        `, prefix+"%").
		Scan(&lastPONumber).Error

	if err != nil {
		log.Error("❌ Failed reading last PO number: " + err.Error())
		return "", err
	}

	log.Infof("📌 Last PO Number from DB = %s", lastPONumber)

	// ==== EXTRACT LAST 4 DIGITS ====
	seq := 0
	if lastPONumber != "" && len(lastPONumber) >= len(prefix)+4 {
		lastSeqStr := lastPONumber[len(lastPONumber)-4:]
		lastSeq, _ := strconv.Atoi(lastSeqStr)
		seq = lastSeq
	}

	// ==== INCREMENT ====
	seq++

	// Final PO number
	newPONumber := fmt.Sprintf("%s%04d", prefix, seq)

	log.Infof("🎉 New PO Number Generated = %s", newPONumber)
	return newPONumber, nil
}

func NewCreatePurchaseOrderService(db *gorm.DB, payload PurchaseOrderPayload, roleName string, createdBy int) (map[string]interface{}, error) {
	log := logger.InitLogger()
	log.Info("🛠️ CreatePurchaseOrderService invoked")
	log.Infof("📥 Payload: %+v", payload)
	log.Infof("👤 Created By: %s", roleName)

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	poNumber, err := GeneratePONumber(db, year, month)
	if err != nil {
		log.Error("❌ Failed to generate PO Number: " + err.Error())
		return nil, err
	}

	log.Infof("🧾 Generated PO Number: %s", poNumber)

	createdAt := now.Format("2006-01-02 15:04:05")

	// INSERT PO HEADER
	var poId int
	err = db.Raw(`
		INSERT INTO "PurchaseOrderManagement"."PurchaseOrders"
		(po_number, "supplierId", branchid, "taxEnabled", "taxRate",
		 "paymentFee", "shippingFee", "subTotal", "taxAmount", "roundOff", total,
		 "poYear", "poMonth", status, "createdAt", "createdBy", "isDelete")
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'OPEN', ?, ?, FALSE)
		RETURNING id
	`,
		poNumber, payload.SupplierId, payload.BranchId,
		payload.TaxEnabled, fmt.Sprintf("%.2f", payload.TaxRate),
		fmt.Sprintf("%.2f", payload.PaymentFee),
		fmt.Sprintf("%.2f", payload.ShippingFee),
		fmt.Sprintf("%.2f", payload.Subtotal),
		fmt.Sprintf("%.2f", payload.TaxAmount),
		fmt.Sprintf("%.2f", payload.RoundOff),
		fmt.Sprintf("%.2f", payload.Total),
		fmt.Sprintf("%d", year),
		fmt.Sprintf("%d", month),
		createdAt, roleName,
	).Scan(&poId).Error

	if err != nil {
		log.Error("❌ Failed inserting PO header: " + err.Error())
		return nil, err
	}

	log.Infof("🧾 Purchase Order ID: %d", poId)

	// INSERT ITEMS
	for _, item := range payload.Items {
		db.Exec(`
			INSERT INTO "PurchaseOrderManagement"."PurchaseOrderItems"
			("purchaseOrderId", "categoryId", "subCategoryId", "productDescription",
			 "unitPrice", quantity, "discountPercent", "discountAmount", "lineTotal",
			 "receivedQuantity", "isClosed", "createdAt")
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, FALSE, ?)
		`,
			poId,
			item.CategoryId, item.SubCategoryId, item.ProductDescription,
			fmt.Sprintf("%.2f", item.UnitPrice),
			fmt.Sprintf("%.2f", item.Quantity),
			fmt.Sprintf("%.2f", item.DiscountPercent),
			fmt.Sprintf("%.2f", item.DiscountAmount),
			fmt.Sprintf("%.2f", item.Total),
			createdAt,
		)
	}

	// INSERT AUDIT
	jsonData, _ := json.Marshal(payload)

	db.Exec(`
		INSERT INTO "PurchaseOrderManagement"."PurchaseOrderAudit"
		("purchaseOrderId", "actionType", "actionDetails", "createdAt", "createdBy")
		VALUES (?, 'CREATE', ?, ?, ?)
	`,
		poId, string(jsonData), createdAt, roleName,
	)

	log.Info("📘 Logged Audit trail")

	// LOG TRANSACTION
	transErr := transactionLogger.LogTransaction(
		db, 1, roleName, 2,
		"Purchase Order Created: "+poNumber,
	)
	if transErr != nil {
		log.Error("⚠️ Transaction Log Failed: " + transErr.Error())
	}

	return map[string]interface{}{
		"poId":     poId,
		"poNumber": poNumber,
	}, nil
}

func NewGetAllPurchaseOrdersService(db *gorm.DB) []map[string]interface{} {
	log := logger.InitLogger()
	log.Info("🛠️ GetAllPurchaseOrdersService invoked")

	var list []map[string]interface{}

	db.Raw(`
		SELECT
			po.id,
			po.po_number,
			po."supplierId",
			s."supplierName",
			s."creditedDays",
			po.branchid,
			b."refBranchId",
			b."refBranchCode",

			po."taxEnabled",
			po."taxRate",
			po."taxAmount",
			po."subTotal",
			po.total,
			po.status,
			po."createdAt",
			po."createdBy",

			-- Ordered Quantity From PO Items
			SUM(COALESCE(poi.quantity::numeric, 0)) AS totalOrderedQty,

			-- Received Quantity From GRN (corrected)
			COALESCE(grni.total_received, 0) AS totalReceivedQty,

			-- Fully Closed?
			CASE 
				WHEN SUM(COALESCE(poi.quantity::numeric, 0)) 
					= COALESCE(grni.total_received, 0)
				THEN TRUE 
				ELSE FALSE 
			END AS isFullyClosed

		FROM "PurchaseOrderManagement"."PurchaseOrders" po

		JOIN public."Branches" b 
			ON b."refBranchId" = po.branchid

		JOIN public."Supplier" s 
			ON s."supplierId" = po."supplierId"

		LEFT JOIN "PurchaseOrderManagement"."PurchaseOrderItems" poi
			ON poi."purchaseOrderId" = po.id

		-- Correct GRN Aggregation
		LEFT JOIN (
			SELECT 
				g."purchaseOrderId",
				COUNT(*) AS total_received
			FROM "PurchaseOrderManagement"."PurchaseOrderGRN" g
			JOIN "PurchaseOrderManagement"."PurchaseOrderGRNItems" gi
				ON gi."grnId" = g.id
			WHERE gi."isDelete" = FALSE
			GROUP BY g."purchaseOrderId"
		) grni ON grni."purchaseOrderId" = po.id

		WHERE po."isDelete" = 'false'

		GROUP BY
			po.id,
			po.po_number,
			po."supplierId",
			s."supplierName",
			s."creditedDays",
			po.branchid,
			b."refBranchId",
			b."refBranchCode",
			po."taxEnabled",
			po."taxRate",
			po."taxAmount",
			po."subTotal",
			po.total,
			po.status,
			po."createdAt",
			po."createdBy",
			grni.total_received     -- IMPORTANT new addition

		ORDER BY po.id DESC;


	`).Scan(&list)

	log.Infof("📊 Retrieved %d purchase orders", len(list))

	return list
}

func NewGetSinglePurchaseOrderService(db *gorm.DB, poId int) (map[string]interface{}, error) {
	log := logger.InitLogger()
	log.Infof("🛠️ Fetching PO ID: %d", poId)

	var header map[string]interface{}
	db.Raw(`
		SELECT * FROM "PurchaseOrderManagement"."PurchaseOrders"
		WHERE id = ? AND "isDelete" = FALSE
	`, poId).Scan(&header)

	if header["id"] == nil {
		return nil, fmt.Errorf("PO not found")
	}

	var items []map[string]interface{}
	db.Raw(`
		SELECT * FROM "PurchaseOrderManagement"."PurchaseOrderItems"
		WHERE "purchaseOrderId" = ?
	`, poId).Scan(&items)

	header["items"] = items

	log.Info("✅ PO fetched successfully")

	return header, nil
}

type GRNItem struct {
	SNo           int     `json:"sNo"`
	LineNo        string  `json:"lineNo"`
	RefNo         string  `json:"refNo"`
	Cost          float64 `json:"cost"`
	ProfitPercent float64 `json:"profitPercent"`
	Total         float64 `json:"total"`

	Design struct {
		Id   any    `json:"id"` // can be int or string
		Name string `json:"name"`
	} `json:"design"`

	Pattern struct {
		Id   any    `json:"id"`
		Name string `json:"name"`
	} `json:"pattern"`

	Variant struct {
		Id   any    `json:"id"`
		Name string `json:"name"`
	} `json:"variant"`

	Color struct {
		Id   any    `json:"id"`
		Name string `json:"name"`
	} `json:"color"`

	Size struct {
		Id   any    `json:"id"`
		Name string `json:"name"`
	} `json:"size"`

	ProductId   int    `json:"productId"`
	ProductName string `json:"productName"`
}

type GRNPayload struct {
	PoId       int       `json:"poId"`
	SupplierId int       `json:"supplierId"`
	BranchId   int       `json:"branchId"`
	Items      []GRNItem `json:"items"`
}

func SafeInt(v any) *int {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case int:
		return &val
	case float64:
		i := int(val)
		return &i
	case string:
		if val == "" {
			return nil
		}
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil
		}
		return &i
	default:
		return nil
	}
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func GenerateSKU(db *gorm.DB, year int, month int) (string, error) {
	log := logger.InitLogger()
	log.Info("🔧 Generating SKU...")

	yearCode := string(rune('A' + (year - 2025))) // 2025=A
	monthCode := string(rune('A' + (month - 1)))  // Jan=A

	prefix := fmt.Sprintf("SS%s%s", yearCode, monthCode)

	log.Infof("🔍 SKU Prefix = %s", prefix)

	var lastSKU string
	err := db.Raw(`
        SELECT sku
        FROM "PurchaseOrderManagement"."PurchaseOrderGRNItems"
        WHERE sku LIKE ?
        ORDER BY id DESC
        LIMIT 1
    `, prefix+"%").Scan(&lastSKU).Error

	if err != nil {
		log.Error("❌ Failed reading last SKU: " + err.Error())
		return "", err
	}

	seq := 0
	if lastSKU != "" && len(lastSKU) >= len(prefix)+5 {
		lastSeqStr := lastSKU[len(lastSKU)-5:]
		lastSeq, _ := strconv.Atoi(lastSeqStr)
		seq = lastSeq
	}

	seq++

	newSKU := fmt.Sprintf("%s%05d", prefix, seq)

	log.Infof("🎉 Generated SKU = %s", newSKU)
	return newSKU, nil
}

func NewCreateGRNService(db *gorm.DB, payload GRNPayload) (map[string]interface{}, error) {
	log := logger.InitLogger()
	log.Info("🛠️ NewCreateGRNService invoked")

	now := time.Now().Format("2006-01-02 15:04:05")

	// INSERT GRN HEADER
	var grnId int
	err := db.Raw(`
		INSERT INTO "PurchaseOrderManagement"."PurchaseOrderGRN"
		(
			"purchaseOrderId", "supplierId", "supplierName", branchid,
			"branchCode", "poNumber", "grnDate", "totalReceivedQty", "createdAt"
		)
		SELECT 
			po.id, 
			po."supplierId",
			s."supplierName",
			po.branchid,
			b."refBranchCode",
			po.po_number,
			?, ?, ?
		FROM "PurchaseOrderManagement"."PurchaseOrders" po
		JOIN public."Supplier" s ON s."supplierId" = po."supplierId"
		JOIN public."Branches" b ON b."refBranchId" = po.branchid
		WHERE po.id = ?
		RETURNING id
	`,
		now,                                   // grnDate
		fmt.Sprintf("%d", len(payload.Items)), // totalReceivedQty
		now,                                   // createdAt
		payload.PoId,
	).Scan(&grnId).Error

	if err != nil {
		log.Error("❌ Failed inserting GRN header: " + err.Error())
		return nil, err
	}

	log.Infof("🆔 GRN Created with ID = %d", grnId)

	// INSERT GRN ITEMS
	for _, item := range payload.Items {

		// 👉 Generate SKU per item
		sku, err := GenerateSKU(db, time.Now().Year(), int(time.Now().Month()))
		if err != nil {
			return nil, err
		}

		err = db.Exec(`
    INSERT INTO "PurchaseOrderManagement"."PurchaseOrderGRNItems"
    (
        "grnId", "purchaseOrderId", "supplierId",
        "lineNo", "refNo",
        "productId", "productName",
        "designId", "designName",
        "patternId", "patternName",
        "varientId", "varientName",
        "colorId", "colorName",
        "sizeId", "sizeName",
        cost, "profitPercent", total,
        "createdAt", "createdBy",
        "updatedAt", "updatedBy",
        "productBranchId", "isDelete",
        quantity,
        sku
    )
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`,
			grnId,
			payload.PoId,
			payload.SupplierId,

			item.LineNo,
			item.RefNo,

			item.ProductId,
			item.ProductName,

			SafeInt(item.Design.Id),
			item.Design.Name,

			SafeInt(item.Pattern.Id),
			item.Pattern.Name,

			SafeInt(item.Variant.Id),
			item.Variant.Name,

			SafeInt(item.Color.Id),
			item.Color.Name,

			SafeInt(item.Size.Id),
			item.Size.Name,

			toString(item.Cost),
			toString(item.ProfitPercent),
			toString(item.Total),

			now,     // createdAt
			"admin", // createdBy (or roleName)
			nil,     // updatedAt
			nil,     // updatedBy

			payload.BranchId,
			false, // isDelete

			1,   // quantity (default)
			sku, // newly generated SKU
		).Error

		if err != nil {
			log.Error("❌ Failed inserting GRN item: " + err.Error())
			return nil, err
		}
	}

	return map[string]interface{}{
		"grnId": grnId,
	}, nil
}

func NewGetAllGRNService(db *gorm.DB) []map[string]interface{} {
	var list []map[string]interface{}
	db.Raw(`
		SELECT grn.*, po.po_number
		FROM "PurchaseOrderManagement"."PurchaseOrderGRN" grn
		JOIN "PurchaseOrderManagement"."PurchaseOrders" po
			ON po.id = grn."purchaseOrderId"
		ORDER BY grn.id DESC
	`).Scan(&list)

	return list
}

func NewGetSingleGRNService(db *gorm.DB, grnId int) (map[string]interface{}, error) {
	var header map[string]interface{}
	db.Raw(`
		SELECT * FROM "PurchaseOrderManagement"."PurchaseOrderGRN"
		WHERE id = ?
	`, grnId).Scan(&header)

	if header["id"] == nil {
		return nil, fmt.Errorf("GRN not found")
	}

	var items []map[string]interface{}
	db.Raw(`
		SELECT * FROM "PurchaseOrderManagement"."PurchaseOrderGRNItems"
		WHERE "grnId" = ?
	`, grnId).Scan(&items)

	header["items"] = items
	return header, nil
}
