package accesstoken

import (
	"fmt"
	"os"
	"strings"
	"time"

	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreateToken generates a JWT token for a given user ID and expiration duration.
func CreateToken(id any, roleId any, branchid any) string {
	log := logger.InitLogger()
	log.Info("🔑 Creating JWT token...")

	jwtKey := []byte(os.Getenv("ACCESS_TOKEN"))
	claims := jwt.MapClaims{
		"id":       id,
		"roleId":   roleId,
		"branchId": branchid,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Error(fmt.Sprintf("❌ Error creating token: %v", err))
		return "Invalid Token"
	}

	log.Infof("✅ Token created successfully for userId=%v, roleId=%v, branchId=%v", id, roleId, branchid)
	return tokenString
}

// ValidateJWT validates the JWT token and checks if it is expired.
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	log := logger.InitLogger()
	log.Infof("🔍 Validating JWT token: %s", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("❌ Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("ACCESS_TOKEN")), nil
	})

	if err != nil {
		log.Error(fmt.Sprintf("❌ JWT parsing failed: %v", err))
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			log.Error("❌ Invalid exp type in claims")
			return nil, fmt.Errorf("invalid exp type")
		}

		expTime := time.Unix(int64(expFloat), 0)
		if time.Now().After(expTime) {
			log.Warn(fmt.Sprintf("⚠️ Token expired at %s", expTime.Format(time.RFC3339)))
			return nil, fmt.Errorf("token expired at %s", expTime.Format(time.RFC3339))
		}

		log.Infof("✅ Token valid for user ID: %v", claims["id"])
	} else {
		log.Error("❌ Token invalid or missing claims")
	}

	return token, nil
}

// JWTMiddleware protects routes by validating JWT tokens from the Authorization header.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.InitLogger()
		log.Info("🔐 JWT Middleware invoked")

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Error("❌ Missing Token in request header")
			c.JSON(200, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			log.Info("✂️ Stripping Bearer prefix")
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}
		log.Infof("📜 Final Token String: %s", tokenString)

		// 🔎 Parse token without validating expiration first, just to read claims
		parsedToken, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			log.Infof("📝 Token Claims Extracted [Pre-Validation] -> id=%v, roleId=%v, branchId=%v, time=%s",
				claims["id"], claims["roleId"], claims["branchId"], time.Now().Format(time.RFC3339))
		}

		// ✅ Now validate the JWT token (signature + expiration)
		token, err := ValidateJWT(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "token expired") {
				log.Warn("⚠️ Token Expired")
				c.JSON(200, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
			log.Error(fmt.Sprintf("❌ Invalid Token: %v", err))
			c.JSON(200, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract the claims (user info) and set it in the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Infof("✅ Setting claims in context: id=%v, roleId=%v, branchId=%v",
				claims["id"], claims["roleId"], claims["branchId"])

			c.Set("id", claims["id"])
			c.Set("roleId", claims["roleId"])
			c.Set("branchId", claims["branchId"])
			c.Set("token", tokenString)
		} else {
			log.Warn("⚠️ Token claims missing or invalid")
		}

		// Proceed to the next handler if the token is valid
		c.Next()
		log.Info("➡️ Passed JWT middleware successfully")
	}
}
