package productController

import (
	"net/http"
	"strconv"

	productModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/products/model"
	productService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/products/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"
)

func CreatePOProductController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, roleId, branchId := getUserContext(c)
		if id == nil {
			return
		}

		var product productModel.POProduct
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := productService.CreatePOProduct(dbConn, &product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create PO product"})
			return
		}

		token := accesstoken.CreateToken(id, roleId, branchId)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "PO product created successfully", "token": token})
	}
}

func GetAllPOProductsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, roleId, branchId := getUserContext(c)
		if id == nil {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		products, err := productService.GetAllPOProducts(dbConn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to fetch PO products"})
			return
		}

		token := accesstoken.CreateToken(id, roleId, branchId)
		c.JSON(http.StatusOK, gin.H{"status": true, "data": products, "token": token})
	}
}

func GetPOProductByIdController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, roleId, branchId := getUserContext(c)
		if id == nil {
			return
		}

		poId := c.Param("id")
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		product, err := productService.GetPOProductById(dbConn, poId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "PO product not found"})
			return
		}

		token := accesstoken.CreateToken(id, roleId, branchId)
		c.JSON(http.StatusOK, gin.H{"status": true, "data": product, "token": token})
	}
}

func UpdatePOProductController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, roleId, branchId := getUserContext(c)
		if id == nil {
			return
		}

		var product productModel.POProduct
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := productService.UpdatePOProduct(dbConn, &product)
		if err != nil {
			if err.Error() == "cannot update a deleted product" {
				c.JSON(http.StatusForbidden, gin.H{"status": false, "message": "This product is deleted and cannot be updated"})
				return
			}
			if err.Error() == "product not found" {
				c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "PO Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update PO Product"})
			return
		}

		token := accesstoken.CreateToken(id, roleId, branchId)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "PO product updated successfully", "token": token})
	}
}

func DeletePOProductController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, roleId, branchId := getUserContext(c)
		if id == nil {
			return
		}

		poId := c.Param("id")
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := productService.DeletePOProduct(dbConn, poId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete PO product"})
			return
		}

		token := accesstoken.CreateToken(id, roleId, branchId)
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "PO product deleted successfully", "token": token})
	}
}

// Utility
func getUserContext(c *gin.Context) (interface{}, interface{}, interface{}) {
	idValue, idExists := c.Get("id")
	roleIdValue, roleIdExists := c.Get("roleId")
	branchIdValue, branchIdExists := c.Get("branchId")

	if !idExists || !roleIdExists || !branchIdExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "User ID, RoleID, Branch ID not found in request context.",
		})
		return nil, nil, nil
	}
	return idValue, roleIdValue, branchIdValue
}

type CheckSKURequest struct {
	FromBranchID int    `json:"fromBranchId" binding:"required"`
	SKU          string `json:"sku" binding:"required"`
}

func CheckSKUInBranchController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CheckSKURequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		product, found, branchName, err := productService.GetProductBySKUInBranch(dbConn, req.FromBranchID, req.SKU)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		if found {
			c.JSON(http.StatusOK, gin.H{
				"status":     true,
				"isPresent":  true,
				"data":       product,
				"branchName": branchName,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":     true,
				"isPresent":  false,
				"data":       product, // original product from other branch
				"branchName": branchName,
				"message":    "Product exists but not in the given branch",
			})
		}

	}
}

func GetBranch4ProductsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// Call service
		products, err := productService.GetProductsByBranchID(dbConn, 4) // branchId = 4
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(products) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "No products found for branch 4",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   products,
		})
	}
}

func CreateStockTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload productModel.StockTransferRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		transferID, err := productService.CreateStockTransfer(dbConn, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     true,
			"message":    "Stock transfer created successfully",
			"transferId": transferID,
		})
	}
}

func GetStockTransferByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {

		idStr := c.Param("id")
		transferId, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"status": false, "message": "Invalid ID"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		transfer, items, err := productService.GetStockTransferByID(dbConn, transferId)
		if err != nil {
			c.JSON(404, gin.H{"status": false, "message": "Stock transfer not found"})
			return
		}

		c.JSON(200, gin.H{
			"status":   true,
			"transfer": transfer,
			"items":    items,
		})
	}
}

func GetStockTransfersController() gin.HandlerFunc {
	return func(c *gin.Context) {

		toBranchIdStr := c.Query("toBranchId")

		toBranchId, err := strconv.Atoi(toBranchIdStr)
		if err != nil {
			c.JSON(400, gin.H{"status": false, "message": "Invalid branch ID"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		transfers, err := productService.GetStockTransfers(dbConn, toBranchId)
		if err != nil {
			c.JSON(500, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"status": true,
			"data":   transfers,
		})
	}
}

func GetAllStockTransfersController() gin.HandlerFunc {
	return func(c *gin.Context) {

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		transfers, err := productService.GetAllStockTransfers(dbConn)
		if err != nil {
			c.JSON(500, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"status": true,
			"data":   transfers,
		})
	}
}

func ReceiveStockProductsController() gin.HandlerFunc {
	return func(c *gin.Context) {

		var payload productModel.ReceiveStockProductsRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"status": false, "message": "Invalid request", "error": err.Error()})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		if err := productService.ReceiveProductsService(dbConn, payload); err != nil {
			c.JSON(500, gin.H{"status": false, "message": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": true, "message": "Products received successfully"})
	}
}

func SaveProductImagesController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n\n🖼️ SaveProductImagesController invoked")

		// Extract token claims
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		branchIdValue, branchIdExists := c.Get("branchId")

		log.Infof("🔍 Context -> id=%v | roleId=%v | branchId=%v", idValue, roleIdValue, branchIdValue)

		if !idExists || !roleIdExists || !branchIdExists {
			log.Warn("❌ Missing context data")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found",
			})
			return
		}

		// Request body
		var body struct {
			FileNames []string `json:"fileNames"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			log.Error("❌ Invalid payload: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request body",
			})
			return
		}

		log.Infof("📦 Received %d file names", len(body.FileNames))

		// DB
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		err := productService.SaveProductImagesService(dbConn, body.FileNames, idValue)
		if err != nil {
			log.Error("❌ Service Error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to save image details",
			})
			return
		}

		token := accesstoken.CreateToken(idValue, roleIdValue, branchIdValue)

		log.Info("✅ Product images saved successfully\n\n")

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Image details saved successfully",
			"token":   token,
		})
	}
}

func GetImagesByProductController() gin.HandlerFunc {
	log := logger.InitLogger()

	return func(c *gin.Context) {
		log.Info("\n📥 GetImagesByProductController invoked")

		productInstanceId := c.Param("productInstanceId")
		if productInstanceId == "" {
			log.Warn("❌ Missing productInstanceId")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "productInstanceId is required",
			})
			return
		}

		log.Infof("🔍 Fetching images for product_instance_id: %s", productInstanceId)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		data, err := productService.GetImagesByProductService(dbConn, productInstanceId)
		if err != nil {
			log.Error("❌ Failed to fetch images: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Unable to fetch images",
			})
			return
		}

		log.Infof("✅ Found %d images for product_instance_id: %s\n", len(data), productInstanceId)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   data,
		})
	}
}
