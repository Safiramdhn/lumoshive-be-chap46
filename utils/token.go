package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID int, secret string, expiryDuration time.Duration) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"sub": userID,                                // Subject (e.g., user ID)
		"exp": time.Now().Add(expiryDuration).Unix(), // Expiration time
		"iat": time.Now().Unix(),                     // Issued at time
	}

	// Create the token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
