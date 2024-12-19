package middleware

import (
	"majin/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CSRFMiddleware(secure bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "GET" {
			csrfToken := utils.GenerateToken()
			c.Cookie(&fiber.Cookie{
				Name:     "_csrf",
				Value:    csrfToken,
				Path:     "/",
				HTTPOnly: true,
				SameSite: "Lax",
				Secure:   secure,
				MaxAge:   int(time.Hour.Seconds()),
			})
			c.Locals("_csrf", csrfToken)
		}
		return c.Next()
	}
}
