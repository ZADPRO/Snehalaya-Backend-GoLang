package productService

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	bulkImageUploadService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/bulkImageHandling/service"
	poModuleModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/poModule/model"
	productModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/products/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

func CreatePOProduct(db *gorm.DB, product *productModel.POProduct) error {
	log := logger.InitLogger()

	// Format: SKU-dd-mm-yy-00001
	today := time.Now()
	datePart := today.Format("02-01-06") // dd-mm-yy

	var count int64
	err := db.Table("POProducts").
		Where(`"isDelete" = false`).
		Count(&count).Error
	if err != nil {
		log.Error("Failed to count existing products: " + err.Error())
		return err
	}

	sequenceNumber := fmt.Sprintf("%05d", count+1)
	generatedSKU := fmt.Sprintf("SKU-%s-%s", datePart, sequenceNumber)

	// Assign generated SKU
	product.PoSKU = generatedSKU

	// Check for duplicates (by poHSN + poDescription only, since SKU is new)
	var existing productModel.POProduct
	err = db.Table("POProducts").
		Where(`"poHSN" = ? AND "poDescription" = ? AND "isDelete" = false`,
			product.PoHSN, product.PoDescription).
		First(&existing).Error

	if err == nil {
		log.Error("Duplicate PO Product found")
		return fmt.Errorf("duplicate value found")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("DB error while checking for duplicates: " + err.Error())
		return err
	}

	// Proceed to insert
	product.CreatedAt = today.Format("2006-01-02 15:04:05")
	product.CreatedBy = "Admin"
	product.IsDelete = false

	log.Info("Generated SKU: " + product.PoSKU)

	return db.Table("POProducts").Create(product).Error
}

func GetAllPOProducts(db *gorm.DB) ([]productModel.POProduct, error) {
	var products []productModel.POProduct
	err := db.Table("POProducts").
		Where(`"isDelete" = false`).
		Find(&products).Error
	return products, err
}

func GetPOProductById(db *gorm.DB, id string) (productModel.POProduct, error) {
	var product productModel.POProduct
	err := db.Table("POProducts").
		Where(`"poId" = ? AND "isDelete" = false`, id).
		First(&product).Error
	return product, err
}

func UpdatePOProduct(db *gorm.DB, product *productModel.POProduct) error {
	log := logger.InitLogger()

	// Check if the product exists and is not deleted
	var existing productModel.POProduct
	err := db.Table("POProducts").
		Where(`"poId" = ?`, product.POId).
		First(&existing).Error

	if err != nil {
		log.Error("PO Product not found: " + err.Error())
		return fmt.Errorf("product not found")
	}

	if existing.IsDelete {
		log.Error("Attempted to update a deleted product")
		return fmt.Errorf("cannot update a deleted product")
	}

	// Set update metadata
	product.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	product.UpdatedBy = "Admin"

	if product.PoSKU != "" {
		log.Warn("PoSKU was passed in update payload. Ignoring it to preserve original value.")
	}

	// Perform the update
	return db.Table("POProducts").
		Where(`"poId" = ?`, product.POId).
		Updates(map[string]interface{}{
			"poName":        product.PoName,
			"poDescription": product.PoDescription,
			"poHSN":         product.PoHSN,
			"poQuantity":    product.PoQuantity,
			"poPrice":       product.PoPrice,
			"poDiscPercent": product.PoDiscPercent,
			"poDisc":        product.PoDisc,
			"poTotalPrice":  product.PoTotalPrice,
			"updatedAt":     product.UpdatedAt,
			"updatedBy":     product.UpdatedBy,
		}).Error
}

func DeletePOProduct(db *gorm.DB, id string) error {
	return db.Table("POProducts").
		Where(`"poId" = ?`, id).
		Updates(map[string]interface{}{
			"isDelete":  true,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
			"updatedBy": "Admin",
		}).Error
}

type ProductWithBranch struct {
	productModel.PurchaseOrderProduct
	RefBranchName string `gorm:"column:refBranchName"`
	IsPresent     bool   `json:"isPresent" gorm:"-"`
}

func GetProductBySKUInBranch(db *gorm.DB, branchID int, sku string) (ProductWithBranch, bool, string, error) {
	var product ProductWithBranch
	var productOtherBranch ProductWithBranch

	err := db.Raw(`
        SELECT p.*, b."refBranchName"
        FROM "purchaseOrderMgmt"."PurchaseOrderAcceptedProducts" p
        JOIN public."Branches" b ON p."productBranchId" = b."refBranchId"
        WHERE p."productBranchId" = ?
        AND p."SKU" = ?
        AND p."isDelete" = false
        LIMIT 1
    `, branchID, sku).Scan(&product).Error

	if err == nil && product.SKU != "" {
		product.IsPresent = true
		return product, true, product.RefBranchName, nil
	}

	err2 := db.Raw(`
        SELECT p.*, b."refBranchName"
        FROM "purchaseOrderMgmt"."PurchaseOrderAcceptedProducts" p
        JOIN public."Branches" b ON p."productBranchId" = b."refBranchId"
        WHERE p."SKU" = ?
        AND p."isDelete" = false
        LIMIT 1
    `, sku).Scan(&productOtherBranch).Error

	if err2 != nil || productOtherBranch.SKU == "" {
		return ProductWithBranch{}, false, "", fmt.Errorf("SKU not found in any branch")
	}

	productOtherBranch.IsPresent = false
	return productOtherBranch, false, productOtherBranch.RefBranchName, nil
}

type Product4Branch struct {
	productModel.PurchaseOrderProduct
}

func GetProductsByBranchID(db *gorm.DB, branchID int) ([]Product4Branch, error) {
	var products []Product4Branch

	// Fetch products only for this branch
	err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts" p`).
		Select(`p.*`).
		Where(`p."productBranchId" = ? AND p."isDelete" = false AND p.status = 'Active'`, branchID).
		Order(`p.product_instance_id`).
		Scan(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func CreateStockTransfer(db *gorm.DB, payload productModel.StockTransferRequest) (int, error) {

	branchCode := payload.BranchDetails.BranchCode

	var count int64
	if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransfers"`).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count stock transfers: %v", err)
	}

	var nextNumber int64 = 10001

	if count > 0 {
		var lastST string
		err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransfers"`).
			Select(`"po_number"`).
			Where(`"po_number" LIKE 'ST-%'`).
			Order(`"po_number" DESC`).
			Limit(1).
			Scan(&lastST).Error

		if err != nil {
			return 0, fmt.Errorf("failed to fetch last stock transfer number: %v", err)
		}

		if lastST != "" {
			parts := strings.Split(lastST, "-")
			if len(parts) >= 3 {
				if num, convErr := strconv.ParseInt(parts[2], 10, 64); convErr == nil {
					nextNumber = num + 1
				}
			}
		}
	}

	poNumber := fmt.Sprintf("ST-%s-%05d", branchCode, nextNumber)

	transfer := productModel.StockTransfer{
		FromBranchID:      payload.BranchDetails.BranchId,
		FromBranchName:    payload.BranchDetails.BranchName,
		FromBranchEmail:   payload.BranchDetails.BranchEmail,
		FromBranchAddress: payload.BranchDetails.BranchAddress,
		ToBranchID:        payload.ReceivedBranchDetails.SupplierId,
		ToBranchName:      payload.ReceivedBranchDetails.SupplierName,
		ToBranchEmail:     payload.ReceivedBranchDetails.SupplierCompanyName,
		ToBranchAddress:   payload.ReceivedBranchDetails.SupplierGSTNumber,
		ModeOfTransport:   payload.TotalSummary.ModeOfTransport,
		SubTotal:          payload.TotalSummary.SubTotal,
		DiscountOverall:   payload.TotalSummary.DiscountOverall,
		TotalAmount:       payload.TotalSummary.TotalAmount,
		PaymentPending:    payload.TotalSummary.PaymentPending,
		PoNumber:          poNumber,
		Status:            payload.TotalSummary.Status,
		CreatedAt:         payload.TotalSummary.CreatedAt,
		CreatedBy:         payload.TotalSummary.CreatedBy,
		UpdatedAt:         payload.TotalSummary.UpdatedAt,
		UpdatedBy:         payload.TotalSummary.UpdatedBy,
		IsDelete:          false,
	}

	if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransfers"`).
		Create(&transfer).Error; err != nil {
		return 0, fmt.Errorf("failed to insert stock transfer: %v", err)
	}

	for _, p := range payload.ProductDetails {

		item := productModel.StockTransferItem{
			StockTransferID:   transfer.StockTransferID,
			ProductInstanceID: 0,
			ProductName:       p.ProductName,
			SKU:               p.SKU,
			IsReceived:        p.IsReceived,
			AcceptanceStatus:  "In Transit",
		}

		if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransferItems"`).
			Create(&item).Error; err != nil {

			return 0, fmt.Errorf("failed to insert item: %v", err)
		}
	}

	return transfer.StockTransferID, nil
}

func GetStockTransferByID(db *gorm.DB, transferId int) (*productModel.StockTransfer, []productModel.StockTransferItem, error) {

	var transfer productModel.StockTransfer
	var items []productModel.StockTransferItem

	// Fetch parent
	if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransfers"`).
		Where(`stock_transfer_id = ? AND is_delete = false`, transferId).
		First(&transfer).Error; err != nil {
		return nil, nil, err
	}

	// Fetch items
	if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransferItems"`).
		Where(`stock_transfer_id = ?`, transferId).
		Find(&items).Error; err != nil {
		return &transfer, nil, err
	}

	return &transfer, items, nil
}

func GetStockTransfers(db *gorm.DB, toBranchId int) ([]productModel.StockTransfer, error) {

	var transfers []productModel.StockTransfer

	query := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransfers"`).
		Where(`"to_branch_id" = ? AND "is_delete" = false`, toBranchId).
		Order(`stock_transfer_id DESC`)

	if err := query.Find(&transfers).Error; err != nil {
		return nil, err
	}

	return transfers, nil
}

func GetAllStockTransfers(db *gorm.DB) ([]productModel.StockTransfer, error) {

	var transfers []productModel.StockTransfer

	// Fetch all parent transfers
	if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransfers"`).
		Where(`is_delete = false`).
		Order(`stock_transfer_id DESC`).
		Find(&transfers).Error; err != nil {
		return nil, err
	}

	// Fetch all items
	var items []productModel.StockTransferItem

	if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransferItems"`).
		Find(&items).Error; err != nil {
		return nil, err
	}

	// Map items under each parent
	itemMap := make(map[int][]productModel.StockTransferItem)
	for _, item := range items {
		itemMap[item.StockTransferID] = append(itemMap[item.StockTransferID], item)
	}

	// Attach items to parent
	for i := range transfers {
		id := transfers[i].StockTransferID
		transfers[i].Items = itemMap[id]
	}

	return transfers, nil
}

func ReceiveProductsService(db *gorm.DB, payload productModel.ReceiveStockProductsRequest) error {

	// 1. Fetch Stock Transfer
	var transfer productModel.StockTransfer
	if err := db.Table(`"purchaseOrderMgmt"."Inventory_StockTransfers"`).
		Where(`stock_transfer_id = ?`, payload.StockTransferId).
		First(&transfer).Error; err != nil {

		return fmt.Errorf("invalid stock transfer ID: %v", err)
	}

	toBranchId := transfer.ToBranchID
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 2. Start DB transaction
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, p := range payload.AllProducts {

		// 3. Update Inventory Stock Transfer Item
		if err := tx.Table(`"purchaseOrderMgmt"."Inventory_StockTransferItems"`).
			Where(`stock_transfer_item_id = ?`, p.StockTransferItemID).
			Updates(map[string]interface{}{
				"is_received":       true,
				"acceptance_status": "Received",
			}).Error; err != nil {

			tx.Rollback()
			return fmt.Errorf("failed to update stock transfer item: %v", err)
		}

		// 4. Update PurchaseOrderAcceptedProducts based on SKU
		if err := tx.Table(`"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts"`).
			Where(`"SKU" = ?`, p.SKU).
			Updates(map[string]interface{}{
				"productBranchId": toBranchId,
				"updatedAt":       currentTime,
				"updatedBy":       "Admin",
			}).Error; err != nil {

			tx.Rollback()
			return fmt.Errorf("failed to update accepted product: %v", err)
		}
	}

	// 5. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func SaveProductImagesService(db *gorm.DB, fileNames []string, createdBy interface{}) error {
	log := logger.InitLogger()
	log.Info("\n🛠️ SaveProductImagesService invoked")

	// Correct SKU extraction
	skuRegex := regexp.MustCompile(`^([A-Za-z0-9-]*\d)`)

	for _, fileName := range fileNames {
		log.Infof("📄 Processing file: %s", fileName)

		upper := strings.ToUpper(fileName)

		match := skuRegex.FindStringSubmatch(upper)
		if len(match) < 2 {
			log.Warn("⚠️ No valid SKU found in: " + fileName)
			continue
		}

		baseSku := match[1] // SS-1125-10487
		log.Infof("🔍 Extracted Base SKU: %s", baseSku)

		var productInstanceID int
		err := db.Table(`"purchaseOrderMgmt"."PurchaseOrderAcceptedProducts"`).
			Select("product_instance_id").
			Where(`"SKU" = ?`, baseSku).
			Scan(&productInstanceID).Error

		if err != nil || productInstanceID == 0 {
			log.Warn("⚠️ No DB match for: " + baseSku)
		} else {
			log.Infof("✅ Matched product_instance_id: %d", productInstanceID)
		}

		// Save DB record
		record := productModel.ProductImage{
			FileName:          &fileName,
			ProductInstanceID: &productInstanceID,
			SkuFound:          &baseSku,
			ExtractedSku:      &baseSku,
			CreatedAt:         func() *string { t := time.Now().Format("2006-01-02 15:04:05"); return &t }(),
			CreatedBy:         func() *string { v := fmt.Sprintf("%v", createdBy); return &v }(),
			IsDelete:          func() *bool { b := false; return &b }(),
		}

		log.Infof("💾 Saving record: %+v", record)

		if err := db.Table(`"purchaseOrderMgmt"."ProductImages"`).Create(&record).Error; err != nil {
			log.Error("❌ Failed saving image details: " + err.Error())
			return err
		}
	}

	log.Info("🎉 All images processed successfully")
	return nil
}

func GetImagesByProductService(db *gorm.DB, productInstanceId string) ([]productModel.ProductImage, error) {
	log := logger.InitLogger()

	var images []productModel.ProductImage

	log.Infof("🔍 Querying images for productInstanceId=%s", productInstanceId)

	err := db.Table(`"purchaseOrderMgmt"."ProductImages"`).
		Where(`product_instance_id = ? AND (is_delete = false OR is_delete IS NULL)`, productInstanceId).
		Order("image_id DESC").
		Scan(&images).Error

	if err != nil {
		log.Error("❌ DB Error: " + err.Error())
		return nil, err
	}

	log.Infof("📸 %d images retrieved", len(images))

	return images, nil
}

type SingleProductWithImages struct {
	poModuleModel.PurchaseOrderAcceptedProductResponse
	Images []ImageResponse `json:"images"`
}

type ImageResponse struct {
	FileName string `json:"fileName"`
	ViewURL  string `json:"viewUrl"`
}

func GetSinglePurchaseOrderAcceptedProductService(db *gorm.DB, productInstanceId int) (*SingleProductWithImages, error) {
	log := logger.InitLogger()
	log.Infof("🔍 Fetching productInstanceId: %d", productInstanceId)

	// ---------- FETCH MAIN PRODUCT ----------
	var product poModuleModel.PurchaseOrderAcceptedProductResponse

	query := `
		SELECT
			ap.product_instance_id,
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
			ap.sub_category_id,
			ap.status,
			ap."createdAt" AS created_at,
			ap."createdBy" AS created_by,
			ap."updatedAt" AS updated_at,
			ap."updatedBy" AS updated_by,
			ap."isDelete" AS is_delete,
			ap.product_name,
			ap."purchaseOrderId",
			ap."SKU",
			ap."productBranchId",
			ap.quantity,
			po."invoiceFinalNumber",
			c."categoryName",
			sc."subCategoryName",
			b."refBranchName" AS branch_name
		FROM "purchaseOrderMgmt"."PurchaseOrderAcceptedProducts" ap
		LEFT JOIN "purchaseOrderMgmt"."PurchaseOrders" po ON ap."purchaseOrderId" = po.purchase_order_id
		LEFT JOIN public."Categories" c ON c."refCategoryid" = ap.category_id
		LEFT JOIN public."SubCategories" sc ON sc."refSubCategoryId" = ap.sub_category_id
		LEFT JOIN public."Branches" b ON b."refBranchId" = ap."productBranchId"
		WHERE ap.product_instance_id = ?
		AND ap."isDelete" = false
	`

	err := db.Raw(query, productInstanceId).Scan(&product).Error
	if err != nil || product.ProductInstanceID == 0 {
		return nil, fmt.Errorf("Product not found")
	}

	// ---------- FETCH IMAGES ----------
	var images []struct {
		FileName string `gorm:"column:file_name"`
	}

	err = db.Table(`"purchaseOrderMgmt"."ProductImages"`).
		Select("file_name").
		Where("product_instance_id = ? AND is_delete = false", productInstanceId).
		Scan(&images).Error

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch images")
	}

	// ---------- GENERATE VIEW URL FOR EACH ----------
	imageResponses := []ImageResponse{}
	for _, img := range images {

		// MinIO object name format
		objectName := "bulk-images/" + img.FileName

		viewURL, err := bulkImageUploadService.GetImageViewURL(objectName, 30)
		if err != nil {
			log.Errorf("⚠️ Failed generating view URL for %s", img.FileName)
			continue
		}

		imageResponses = append(imageResponses, ImageResponse{
			FileName: img.FileName,
			ViewURL:  viewURL,
		})
	}

	// Return combined response
	return &SingleProductWithImages{
		PurchaseOrderAcceptedProductResponse: product,
		Images:                               imageResponses,
	}, nil
}
