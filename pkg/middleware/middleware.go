package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kiranpawar037/task-managements-service/pkg/config"
	"github.com/kiranpawar037/task-managements-service/pkg/models"
	"gorm.io/gorm"
)

func JWTMiddlewareUser(db *gorm.DB) gin.HandlerFunc {
	// Load JWT secret
	cfg, err := config.Env()
	if err != nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to load config: %v", err),
			})
			c.Abort()
		}
	}

	jwtSecret := []byte(cfg.JWT.Secret)

	return func(c *gin.Context) {

		// Get Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}

		// Remove Bearer prefix
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Extract email
		email, ok := claims["email"].(string)
		if !ok || email == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Email missing in token"})
			return
		}

		// Load user from DB
		var user models.User
		if err := db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			}
			return
		}

		// Save user in context
		c.Set("user", &user)

		c.Next()
	}
}
