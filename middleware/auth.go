package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			return c.Next()
		}

		if email, ok := session.Get("email").(string); ok {
			c.Locals("email", email)
			c.Locals("authenticated", true)
		} else {
			c.Locals("authenticated", false)
		}

		return c.Next()
	}
}
