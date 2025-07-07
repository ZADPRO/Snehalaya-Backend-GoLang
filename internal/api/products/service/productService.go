package productService

import (
	"fmt"
	"time"

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
