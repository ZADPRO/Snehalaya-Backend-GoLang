package controller

import (
	"net/http"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/service"
)

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	admin := service.GetAdminService()
	w.Write([]byte("Admin Name" + admin.Name))
}
