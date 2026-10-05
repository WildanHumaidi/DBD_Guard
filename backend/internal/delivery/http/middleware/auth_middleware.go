package middleware

import (
	"fmt"
	"strings"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization header format"})
			}
			tokenString = parts[1]
		} else if qToken := c.Query("token"); qToken != "" {
			tokenString = qToken
		} else if cToken := c.Cookies("token"); cToken != "" {
			tokenString = cToken
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.Env.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Parse user_id string from claims into uuid.UUID
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id in token"})
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Malformed user_id in token"})
		}

		// Store typed values in context for next handlers
		c.Locals("user_id", userID)
		c.Locals("role", claims["role"])

		if locIDStr, ok := claims["location_id"].(string); ok {
			if locID, err := uuid.Parse(locIDStr); err == nil {
				c.Locals("location_id", locID)
			}
		}

		return c.Next()
	}
}

// RequireRole checks if the user has one of the allowed roles
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Role not found in context"})
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied: insufficient permissions"})
		}

		return c.Next()
	}
}

// GetUserID is a helper to safely extract user_id from fiber context locals.
// Should only be called inside routes protected by the Protected() middleware.
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}

// GetLocationID is a helper to safely extract location_id from fiber context locals.
func GetLocationID(c *fiber.Ctx) (*uuid.UUID, error) {
	locID, ok := c.Locals("location_id").(uuid.UUID)
	if !ok {
		return nil, nil // Not all users have a location_id
	}
	return &locID, nil
}

// GetUserRole is a helper to safely extract role from fiber context locals.
func GetUserRole(c *fiber.Ctx) string {
	role, ok := c.Locals("role").(string)
	if !ok {
		return ""
	}
	return role
}
