package service

import "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/repository"

type Admin struct {
	Name string
}

func GetAdminService() Admin {
	name := repository.GetAdminFromDb()
	return Admin{Name: name}
}
