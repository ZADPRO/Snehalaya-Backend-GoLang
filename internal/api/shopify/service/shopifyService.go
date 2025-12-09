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

	// ✅ Step 1: Create the product
	createdProduct, err := client.Product.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	log.Printf("✅ Product created successfully! Id: %d, Title: %s\n", createdProduct.Id, createdProduct.Title)

	// ✅ Step 2: Get your Shopify location ID
	locations, err := client.Location.List(ctx, nil)
	if err != nil {
		return createdProduct, fmt.Errorf("failed to get locations: %v", err)
	}
	if len(locations) == 0 {
		return createdProduct, fmt.Errorf("no Shopify locations found")
	}
	locationID := locations[0].Id // 🟢 use .Id instead of .ID
	log.Printf("📍 Using location ID: %d\n", locationID)

	for _, variant := range createdProduct.Variants {
		if variant.InventoryItemId == 0 {
			log.Printf("⚠️ Variant %d has no InventoryItemId, skipping...\n", variant.Id)
			continue
		}

		// Enable tracking
		tracked := true
		invItem := goshopify.InventoryItem{
			Id:      variant.InventoryItemId,
			Tracked: &tracked,
		}
		_, err = client.InventoryItem.Update(ctx, invItem)
		if err != nil {
			log.Printf("⚠️ Failed to enable inventory tracking for variant %d: %v\n", variant.Id, err)
			continue
		}

		// Connect inventory to location
		connectReq := goshopify.InventoryLevel{
			InventoryItemId: variant.InventoryItemId,
			LocationId:      locationID,
		}
		_, err = client.InventoryLevel.Connect(ctx, connectReq)
		if err != nil {
			log.Printf("⚠️ Failed to connect inventory for variant %d: %v\n", variant.Id, err)
			continue
		}

		log.Printf("\n\n\n\n\n\nInventory Qnty", variant.InventoryQuantity)

		setReq := goshopify.InventoryLevel{
			InventoryItemId: variant.InventoryItemId,
			LocationId:      locationID,
			Available:       1,
		}

		_, err = client.InventoryLevel.Set(ctx, setReq)
		if err != nil {
			log.Printf("⚠️ Failed to set inventory for variant %d: %v\n", variant.Id, err)
			continue
		}

		log.Printf("✅ Inventory SET to %d for variant %d\n", variant.InventoryQuantity, variant.Id)

		if err != nil {
			log.Printf("⚠️ Failed to adjust inventory for variant %d: %v\n", variant.Id, err)
			continue
		}

		log.Printf("✅ Inventory set to %d for variant %d\n", variant.InventoryQuantity, variant.Id)
	}

	log.Println("🎉 Product created with tracked inventory!")
	return createdProduct, nil
}
