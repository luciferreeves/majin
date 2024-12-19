package router

import (
	"majin/internal/auth"
	"majin/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Initialize(router *fiber.App) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/mail")
	})

	router.Get("/login", func(c *fiber.Ctx) error {
		return templates.Login().Render(c.Context(), c.Response().BodyWriter())
	})

	authRouter := router.Group("/auth")
	authRouter.Post("/login", func(c *fiber.Ctx) error {
		// Retrieve the session store from locals
		store, ok := c.Locals("session").(*session.Store)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Session initialization error")
		}

		// Create a new session
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Session creation error")
		}

		email := c.FormValue("email")
		password := c.FormValue("password")

		if email == "me@shi.foo" && password == "me@shi.foo" {
			// Set session values
			sess.Set("email", email)

			// Save the session
			if err := sess.Save(); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Session save error")
			}

			// If it's an HTMX request, just confirm success
			if c.Get("HX-Request") == "true" {
				return c.SendStatus(fiber.StatusOK)
			}

			// For non-HTMX requests, redirect
			return c.Redirect("/mail")
		}

		// Return unauthorized status with error message
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	})

	authRouter.Post("/logout", func(c *fiber.Ctx) error {
		// Retrieve the session store from locals
		store, ok := c.Locals("session").(*session.Store)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Session initialization error")
		}

		// Get the current session
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Session retrieval error")
		}

		// Destroy the session
		if err := sess.Destroy(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Logout failed")
		}

		// If it's an HTMX request, render the login template
		if c.Get("HX-Request") == "true" {
			return templates.Login().Render(c.Context(), c.Response().BodyWriter())
		}

		// For non-HTMX requests, redirect to login
		return c.Redirect("/login")
	})

	mail := router.Group("/mail", auth.RequireAuth)

	mail.Get("/", func(c *fiber.Ctx) error {
		return templates.Mail().Render(c.Context(), c.Response().BodyWriter())
	})

	router.Use(func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return templates.NotFound().Render(c.Context(), c)
	})
}
