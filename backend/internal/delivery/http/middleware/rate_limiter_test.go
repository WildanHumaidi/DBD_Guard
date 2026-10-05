package middleware_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPublicRateLimiter(t *testing.T) {
	app := fiber.New()

	app.Use(middleware.PublicRateLimiter())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	// Normal request -> success
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1") // Mock IP

	// Burst request (max is 20)
	for i := 0; i < 20; i++ {
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "success", string(body))
	}

	// 21st request should be rate limited (429)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)

	// Another IP should still work
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-Forwarded-For", "192.168.1.2")
	resp2, err := app.Test(req2, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	// After window -> accepted (this is tricky to test without waiting 1 minute or mocking time,
	// but Fiber's limiter uses sliding windows or fixed windows with real time.
	// We can test the error message structure.
	
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Terlalu banyak permintaan")
}
