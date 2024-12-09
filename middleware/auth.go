package middleware

import (
	"fmt"
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

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			a.respondWithError(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Get user ID from claims
		userID, ok := claims["sub"].(string)
		if !ok {
			a.respondWithError(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Fetch user data from Redis
		userData, err := a.cache.Get(userID)
		if err != nil {
			a.logger.Error("Failed to get user data from Redis", zap.Error(err))
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
