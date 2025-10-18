package config

import (
	"log"

	goshopify "github.com/bold-commerce/go-shopify/v4"

)

var ShopifyClient *goshopify.Client

// Initialize Shopify Client using private app credentials
func InitShopifyClient() {
	app := goshopify.App{
		ApiKey:   "your_api_key_here",
		Password: "your_admin_api_token_here",
	}

	shopName := "yourstore.myshopify.com"

	client, err := goshopify.NewClient(app, shopName, "")
	if err != nil {
		log.Fatalf("❌ Failed to create Shopify client: %v", err)
	}

	ShopifyClient = client
	log.Println("✅ Shopify client initialized successfully")
}
