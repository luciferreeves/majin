package auth

import "github.com/gofiber/fiber/v2"

func RequireAuth(c *fiber.Ctx) error {
	if !IsAuthenticated(c) {
		if c.Get("HX-Request") == "true" {
			c.Set("HX-Redirect", "/login")
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Redirect("/login")
	}
	return c.Next()
}

func IsAuthenticated(c *fiber.Ctx) bool {
	auth, ok := c.Locals("authenticated").(bool)
	return ok && auth
}

func GetUserEmail(c *fiber.Ctx) string {
	email, _ := c.Locals("email").(string)
	return email
}
