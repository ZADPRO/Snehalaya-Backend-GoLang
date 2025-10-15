package productController

import (
	"net/http"

	productModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/products/model"
	productService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/products/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
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
