package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang-chap46/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	logger *zap.Logger
	cache  database.Cacher
	secret string
}

func NewAuthMiddleware(logger *zap.Logger, cache database.Cacher, secret string) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
		cache:  cache,
		secret: secret,
	}
}

func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			a.respondWithError(c, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Extract token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			a.respondWithError(c, http.StatusUnauthorized, "Bearer token is required")
			return
		}

		// Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is valid
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.secret), nil
		})
		if err != nil || !token.Valid {
			a.logger.Error("Invalid token", zap.Error(err))
			a.respondWithError(c, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		log.Printf("token: %v", token)
		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			a.logger.Error("Invalid token claims", zap.Any("claims", token.Claims))
			a.respondWithError(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Safely extract userID
		userID, ok := claims["sub"].(string)
		if !ok {
			subFloat, ok := claims["sub"].(float64)
			if !ok {
				a.respondWithError(c, http.StatusUnauthorized, "Invalid sub claim")
				return
			}
			userID = fmt.Sprintf("%.0f", subFloat) // Convert to string
		}

		// Log and fetch from Redis
		log.Printf("UserID: %s", userID)
		userData, err := a.cache.Get(userID)
		if err != nil {
			a.logger.Error("Failed to fetch user data from cache", zap.Error(err))
			a.respondWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Add user data to context
		c.Set("User", userData)

		// Proceed to next handler
		c.Next()
	}
}

func (a *AuthMiddleware) respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
