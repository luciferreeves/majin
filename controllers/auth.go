package controllers

import (
	"majin/internal/auth"
	"majin/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func LoginView(c *fiber.Ctx) error {
	if auth.IsAuthenticated(c) {
		return c.Redirect("/mail")
	}
	return templates.Login().Render(c.Context(), c.Response().BodyWriter())
}

func Login(c *fiber.Ctx) error {
	store, ok := c.Locals("session").(*session.Store)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Session initialization error")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session creation error")
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "me@shi.foo" && password == "me@shi.foo" {
		sess.Set("email", email)
		sess.Set("password", password)
		c.Locals("authenticated", true)

		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Session save error")
		}

		if c.Get("HX-Request") == "true" {
			return templates.Mail().Render(c.Context(), c.Response().BodyWriter())
		}

		return c.Redirect("/mail")
	}

	return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
}

func Logout(c *fiber.Ctx) error {
	store, ok := c.Locals("session").(*session.Store)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Session initialization error")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session retrieval error")
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Logout failed")
	}

	if c.Get("HX-Request") == "true" {
		return templates.Login().Render(c.Context(), c.Response().BodyWriter())
	}

	return c.Redirect("/login")
}
