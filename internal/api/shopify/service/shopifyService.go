package shopifyService

import (
	"context"
	"fmt"
	"log"

	shopifyConfig "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/shopify"
	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func GetAllProducts() ([]goshopify.Product, error) {
	ctx := context.Background()
	client := shopifyConfig.ShopifyClient
	if client == nil {
		log.Println("⚠️ Shopify client not initialized")
		return nil,
			fmt.Errorf("shopify client not initialized")
	}
	log.Println("✅ Shopify client initialized successfully, fetching products...")
	// Example: Count total products (for verification)
	count, err := client.Product.Count(ctx, nil)
	if err != nil {
		log.Printf("❌ Failed to fetch product count: %v\n", err)
		return nil, err
	}
	log.Printf("📦 Total products available: %d\n", count)
	// Fetch all products
	products, err := client.Product.List(ctx, nil)
	if err != nil {
		log.Printf("❌ Error fetching products: %v\n", err)
		return nil, err
	}
	log.Printf("✅ Successfully fetched %d products\n", len(products))
	return products, nil
}

func CreateProduct(product goshopify.Product) (*goshopify.Product, error) {
	ctx := context.Background()
	client := shopifyConfig.ShopifyClient

	if client == nil {
		return nil, fmt.Errorf("shopify client not initialized")
	}

	createdProduct, err := client.Product.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Product created successfully! Id: %d, Title: %s\n", createdProduct.Id, createdProduct.Title)
	return createdProduct, nil
}
