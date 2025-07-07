package routes

import (
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/controller"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/admin")

	route.POST("/login", controller.AdminLoginController())

	route.POST("/forgot-password", controller.ForgotPasswordController())
	route.POST("/verify-otp", controller.VerifyOtpController())
	route.POST("/reset-password", controller.ResetPasswordController())

}
