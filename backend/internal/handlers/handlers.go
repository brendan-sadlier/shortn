package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Handler struct holds dependencies for handlers
type Handler struct{}

// New creates a new handler instance
func New() *Handler {
	return &Handler{}
}

// HealthCheck is a simple health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// GetUserInfo is a protected endpoint that requires authentication
// It returns comprehensive user information for the dropdown menu
func (h *Handler) GetUserInfo(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userId, exists := c.Get("userId")
	if !exists {
		fmt.Println("Error: User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get the claims from context
	claimsObj, exists := c.Get("claims")
	if !exists {
		fmt.Println("Warning: Claims not found in context, continuing with userID only")
		// If claims aren't available, we can still continue with just the user ID
		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":        userId,
				"email":     "unknown@example.com",
				"name":      "User",
				"avatarUrl": "",
				"role":      "user",
			},
		})
		return
	}

	// Type assertion for claims
	claims, ok := claimsObj.(jwt.MapClaims)
	if !ok {
		fmt.Println("Warning: Failed to assert claims type, continuing with userID only")
		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":        userId,
				"email":     "unknown@example.com",
				"name":      "User",
				"avatarUrl": "",
				"role":      "user",
			},
		})
		return
	}

	// Extract user information from claims
	var (
		email     string
		name      string
		avatarUrl string
		role      string = "user" // Default role
		createdAt string
	)

	// Get email
	if emailClaim, ok := claims["email"].(string); ok {
		email = emailClaim
	}

	// Get name - could be in different claim fields depending on Supabase configuration
	if nameClaim, ok := claims["name"].(string); ok {
		name = nameClaim
	} else if userNameClaim, ok := claims["user_name"].(string); ok {
		name = userNameClaim
	} else if fullNameClaim, ok := claims["full_name"].(string); ok {
		name = fullNameClaim
	} else {
		// Fallback to email name part if no name is available
		if email != "" {
			parts := strings.Split(email, "@")
			if len(parts) > 0 {
				name = parts[0]
			} else {
				name = "User"
			}
		} else {
			name = "User"
		}
	}

	// Get avatar URL
	if avatarClaim, ok := claims["avatar_url"].(string); ok {
		avatarUrl = avatarClaim
	} else if pictureClaim, ok := claims["picture"].(string); ok {
		avatarUrl = pictureClaim
	}

	// Get role/permissions
	if roleClaim, ok := claims["role"].(string); ok {
		role = roleClaim
	} else if appRoleClaim, ok := claims["app_role"].(string); ok {
		role = appRoleClaim
	} else if appMetadata, ok := claims["app_metadata"].(map[string]interface{}); ok {
		if roleFromMeta, ok := appMetadata["role"].(string); ok {
			role = roleFromMeta
		}
	}

	// Get created timestamp if available
	if createdAtClaim, ok := claims["created_at"].(string); ok {
		createdAt = createdAtClaim
	}

	// Return comprehensive user info for dropdown menu
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        userId,
			"email":     email,
			"name":      name,
			"avatarUrl": avatarUrl,
			"role":      role,
			"createdAt": createdAt,
		},
	})
}
