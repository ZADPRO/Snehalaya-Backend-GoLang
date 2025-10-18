package shopifyconfig

import (
	"context"
	"fmt"
	"log"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/joho/godotenv"
)

var (
	ShopifyClient *goshopify.Client
	ShopName      string
	APIToken      string
	APISecretKey  string
)

func Init() {
	_ = godotenv.Load()

	ShopName = os.Getenv("SHOPIFY_SHOP_NAME")
	fmt.Println("\n\nShop Name : ", ShopName)
	APIToken = os.Getenv("SHOPIFY_API_TOKEN")
	fmt.Println("\n\nShop Name : ", APIToken)
	APISecretKey = os.Getenv("SHOPIFY_API_KEY")
	fmt.Println("\n\nShop Name : ", APISecretKey)

	if ShopName == "" || APIToken == "" {
		log.Fatal("Missing SHOPIFY_SHOP_NAME or SHOPIFY_API_TOKEN in environment")
	}

	fmt.Println("API Token from .env:", APIToken)

	app := goshopify.App{
		ApiKey:   APISecretKey,
		Password: "",
	}
	fmt.Println("API Token from .env:", APIToken)

	client, err := goshopify.NewClient(app, ShopName, APIToken)

	if err != nil {
		log.Fatalf("failed to create shopify client: %v", err)
	}

	ShopifyClient = client

	// optional test
	ctx := context.Background()
	if _, err := client.Product.Count(ctx, nil); err != nil {
		log.Fatalf("shopify connectivity test failed: %v", err)
	}
}
