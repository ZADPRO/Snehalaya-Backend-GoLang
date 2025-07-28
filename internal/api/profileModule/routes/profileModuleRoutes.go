package profileModuleRoutes

import (
	accesstoken "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/AccessToken"
	"github.com/gin-gonic/gin"
)

func ProfileModuleRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin/profile")

	route.POST("details", accesstoken.JWTMiddleware())
}
