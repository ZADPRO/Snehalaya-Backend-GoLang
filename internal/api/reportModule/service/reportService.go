package reportService

import (
	reportModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/model"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"gorm.io/gorm"
)

// PRODUCT REPORT SERVICE
func GetAllProductReportsService(db *gorm.DB, productReport *reportModel.ProductsReportPayload, roleName string) (map[string]interface{}, error) {
	log := logger.InitLogger()
	log.Info("Fetching all product reports with pagination")

	var totalCount int64
	var resultDataOfProducts []reportModel.PurchaseOrderResponse

	countQuery := `
		SELECT COUNT(*) 
		FROM "purchaseOrder"."ProductsDummyAcceptance" pda
		WHERE pda."acceptanceStatus" = 'Received'
	`

	if err := db.Raw(countQuery).Scan(&totalCount).Error; err != nil {
		return nil, err
	}

	// OFFSET AND LIMIT
	offset := (productReport.PaginationOffset - 1) * productReport.PaginationLimit
	limit := productReport.PaginationLimit

	dataQuery := `
	SELECT * 
		FROM "purchaseOrder"."ProductsDummyAcceptance" pda
		WHERE pda."acceptanceStatus" = 'Received'
		ORDER BY pda."dummyProductsId" ASC
		LIMIT ? OFFSET ?
	`

	if err := db.Raw(dataQuery, limit, offset).Scan(&resultDataOfProducts).Error; err != nil {
		return nil, err
	}

	log.Infof("Total count of products: %d", totalCount)
	log.Infof("Fetched %d products with offset %d and limit %d", len(resultDataOfProducts), offset, limit)
	log.Infof("Result of the products: %+v", resultDataOfProducts)

	return map[string]interface{}{
		"totalCount": totalCount,
		"data":       resultDataOfProducts,
		"page":       productReport.PaginationOffset,
		"limit":      productReport.PaginationLimit,
	}, nil
}
