package routes

import (
	"net/http"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/controller"
)

func RegisterAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin", controller.GetAdmin)
}
