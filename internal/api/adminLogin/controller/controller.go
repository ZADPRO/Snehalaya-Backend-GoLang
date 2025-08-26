package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/model"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/adminLogin/service"
	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/db"
	becrypt "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Bcrypt"
	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	mailService "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/MailService"
	"github.com/gin-gonic/gin"

)

func AdminLoginController() gin.HandlerFunc {

	log := logger.InitLogger()
	return func(c *gin.Context) {
		var reqVal model.AdminLoginReq

		log.Info("\n\nAdmin Login Controller -> \n================")
		// ERROR HANDLING - STATUS CODE IN PARAMS
		if err := c.BindJSON(&reqVal); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Something went wrong, Try again ... " + err.Error(),
			})
			return
		}

		dbConnt, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.AdminLoginService(dbConnt, reqVal)
		log.Info("Response for controller -> ", resVal)

		response := gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		if resVal.Status {
			response["user"] = resVal.User
			response["token"] = resVal.Token
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// Send OTP (Forgot Password Step 1)

func ForgotPasswordController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid email format"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var user model.AdminLoginModelReq
		err := dbConn.Raw(`SELECT u."refUserId", u."refUserStatus", ucd."refUCDEmail"
				FROM "Users" u
				JOIN "refUserCommunicationDetails" ucd ON u."refUserId" = ucd."refUserId"
				WHERE ucd."refUCDEmail" = ? AND u."refUserStatus" = 'Active'
				LIMIT 1;
				`, req.Email).
			Scan(&user).Error

		if err != nil || user.UserId == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Email not found or inactive"})
			return
		}

		otp := fmt.Sprintf("%04d", rand.Intn(10000))
		expiry := time.Now().Add(2 * time.Minute).Format("2006-01-02 15:04:05")

		otpEntry := model.OTPVerification{
			Email:      req.Email,
			OTP:        otp,
			ExpiresAt:  expiry,
			IsVerified: false,
			CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:  req.Email,
		}
		dbConn.Create(&otpEntry)

		html := fmt.Sprintf(`
		<table width="100%%" cellspacing="0" cellpadding="0" style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 5px;">
			<tr>
				<td align="center">
					<table width="600" cellspacing="0" cellpadding="0" style="background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
						<tr>
							<td style="background-color: #8B0000; padding: 20px; text-align: center; color: white;">
								<h2 style="margin: 0;">🔑 Password Reset Request</h2>
							</td>
						</tr>
						<tr>
							<td style="padding: 30px; color: #333;">
								<p>Dear <strong>User</strong>,</p>

								<p>We received a request to reset your password for your <strong>Snehalayaa Silks ERP</strong> account.</p>

								<p style="margin: 20px 0; font-size: 16px; text-align: center;">
									Your One-Time Password (OTP) is:
								</p>

								<p style="text-align: center; font-size: 24px; font-weight: bold; color: #8B0000; margin: 10px 0;">
									%s
								</p>

								<p style="text-align: center; color: #555;">This OTP will expire in <strong>2 minutes</strong>.</p>

								<p>If you did not request a password reset, please ignore this email. Your account will remain secure.</p>

								<hr style="margin: 30px 0; border: none; border-top: 1px solid #ccc;" />

								<p style="font-style: italic; color: #555; text-align: center;">
									Your security is our priority at Snehalayaa Silks.
								</p>

								<p style="margin-top: 30px;">Warm regards,</p>
								<p><strong>Support Team</strong><br/>Snehalayaa Silks ERP Project</p>
							</td>
						</tr>
						<tr>
							<td style="background-color: #f2f2f2; text-align: center; padding: 15px; font-size: 12px; color: #999;">
								© 2025 Snehalayaa Silks. All rights reserved.
							</td>
						</tr>
					</table>
				</td>
			</tr>
		</table>
	`, otp)
		if !mailService.MailService(req.Email, html, "Password Reset OTP") {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to send OTP email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "OTP sent to your email"})
	}
}

// Verify OTP

func VerifyOtpController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid data"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var otpRecord model.OTPVerification
		err := dbConn.
			Where("email = ? AND otp = ? AND is_verified = false", req.Email, req.OTP).
			Order("expires_at desc").First(&otpRecord).Error

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid OTP"})
			return
		}

		expiryTime, _ := time.Parse("2006-01-02 15:04:05", otpRecord.ExpiresAt)
		if time.Now().After(expiryTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "OTP expired"})
			return
		}

		otpRecord.IsVerified = true
		otpRecord.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		otpRecord.UpdatedBy = req.Email
		dbConn.Save(&otpRecord)

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "OTP verified successfully"})
	}
}

//  Reset Password

func ResetPasswordController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email       string `json:"email"`
			NewPassword string `json:"newPassword"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid input"})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var verifiedOtp model.OTPVerification
		err := dbConn.Where("email = ? AND is_verified = true", req.Email).
			Order("expires_at desc").First(&verifiedOtp).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "OTP not verified"})
			return
		}

		hashedPassword, err := becrypt.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error hashing password"})
			return
		}

		err = dbConn.Exec(`
			WITH target_user AS (
				SELECT u."refUserId"
				FROM "Users" u
				JOIN "refUserCommunicationDetails" comm ON u."refUserId" = comm."refUserId"
				WHERE comm."refUCDEmail" = ?
				AND (u."isDelete" IS FALSE OR u."isDelete" IS NULL)
				ORDER BY u."refUserId"
				LIMIT 1
			)
			UPDATE "refUserAuthCred"
			SET "refUACPassword" = ?, "refUACHashedPassword" = ?
			WHERE "refUserId" = (SELECT "refUserId" FROM target_user)`,
			req.Email, req.NewPassword, hashedPassword).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Password update failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Password reset successfully"})
	}
}
