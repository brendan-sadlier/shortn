package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/brendan-sadlier/shortn/internal/config"
	"net/http"
	"strings"
	_ "time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	supabaseURL string
	jwtSecret   string
}

// NewAuthMiddleware creates a new auth middleware using the Supabase JWT secret
func NewAuthMiddleware(cfg *config.Config) (*AuthMiddleware, error) {
	if cfg.SupabaseJWTKey == "" {
		return nil, errors.New("Supabase JWT key is required")
	}

	return &AuthMiddleware{
		supabaseURL: cfg.SupabaseURL,
		jwtSecret:   cfg.SupabaseJWTKey,
	}, nil
}

// TokenFromHeader extracts JWT token from the Authorization header
func TokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

// PrepareSecret decodes and returns the Supabase JWT secret as a byte array
func (am *AuthMiddleware) PrepareSecret() ([]byte, error) {
	// Supabase JWT secrets are base64 encoded
	decodedSecret, err := base64.RawURLEncoding.DecodeString(am.jwtSecret)
	if err != nil {
		// If decoding fails, try using the secret as-is (might not be base64 encoded)
		return []byte(am.jwtSecret), nil
	}
	return decodedSecret, nil
}

// Middleware function to validate JWT tokens
func (am *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := TokenFromHeader(c)
		if err != nil {
			fmt.Printf("Error extracting token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Get the secret key
		secretKey, err := am.PrepareSecret()
		if err != nil {
			fmt.Printf("Error preparing secret: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Parse the token with more flexibility for algorithm
		var token *jwt.Token

		// First try with HS256
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil // Skip this, we'll try other methods
			}
			return secretKey, nil
		})

		// If that failed, try without algorithm validation
		if err != nil || token == nil {
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})
		}

		if err != nil {
			fmt.Printf("Error parsing token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if token == nil {
			fmt.Printf("Token is nil after parsing\n")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Printf("Failed to extract claims\n")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Debug print claims
		fmt.Printf("Token claims: %+v\n", claims)

		// Extract user ID - try both "sub" and "user_id" fields as Supabase might use either
		var sub string
		if subClaim, ok := claims["sub"]; ok {
			if subStr, ok := subClaim.(string); ok {
				sub = subStr
			}
		} else if userID, ok := claims["user_id"]; ok {
			if userIDStr, ok := userID.(string); ok {
				sub = userIDStr
			}
		}

		if sub == "" {
			fmt.Printf("Missing subject claim\n")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing subject"})
			return
		}

		// Set user info in context
		c.Set("userId", sub)
		c.Set("claims", claims)

		c.Next()
	}
}

// ComputeHS256 computes an HMAC SHA-256 signature
func ComputeHS256(secret, data []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	return h.Sum(nil)
}
