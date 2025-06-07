package main

import (
	"log"
	"net/http"

	adminRoutes "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/routes"
)

func main() {
	mux := http.NewServeMux()
	adminRoutes.RegisterAdminRoutes(mux)

	log.Println("🚀 Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
