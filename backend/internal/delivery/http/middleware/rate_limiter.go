package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// PublicRateLimiter returns a rate limiter middleware for public endpoints
func PublicRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,               // Max 20 requests
		Expiration: 1 * time.Minute,  // per 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.Get("X-Forwarded-For")
			if ip != "" {
				return ip
			}
			return c.IP() // Use client IP as key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Terlalu banyak permintaan. Silakan coba lagi dalam beberapa menit.",
			})
		},
	})
}
