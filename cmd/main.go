package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/routes"
	configMinio "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/config"
	posManagementRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/posManagement/routes"
	productRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/products/routes"
	minioConfig "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/productsImageUpload/config"
	imageUploadRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/productsImageUpload/routes"
	profileModuleRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/profileModule/routes"
	purchaseOrderRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/purchaseOrderModule/routes"
	reportRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/routes"
	settingsRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/settingModule/routes"
	supplierRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/supplierModule/routes"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}
	r.SetTrustedProxies(nil)

	// INIT DB
	db.InitDB()

	//MIN IO INIT
	minioConfig.InitMinio()
	configMinio.InitMinio()

	// CORS CONFIG
	r.Use(cors.New(cors.Config{
		// 	AllowOrigins:     []string{"http://localhost:3000"}, // Change to your allowed origin
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// API CALLS
	routes.RegisterAdminRoutes(r)
	settingsRoutes.SettingsAdminRoutes(r)
	supplierRoutes.SupplierRoutes(r)
	imageUploadRoutes.ImageUploadRoutes(r)
	productRoutes.ProductManagementRoutes(r)
	purchaseOrderRoutes.PurhcaseOrderRoutes(r)
	profileModuleRoutes.ProfileModuleRoutes(r)
	posManagementRoutes.POSManagementRoutes(r)
	reportRoutes.ReportRoutes(r)

	// PING PONG API CALL FOR TESTING
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Pong from User Service",
		})
	})

	// RUN SERVER AND LOG MESSAGE
	fmt.Println("Server is Running at Port : " + os.Getenv("PORT"))
	r.Run("0.0.0.0:" + os.Getenv("PORT"))
}
