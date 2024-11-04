package router

import (
	"majin/internal/auth"
	"majin/templates"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {

	mail := router.Group("/mail", auth.RequireAuth)

	mail.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("You are viewing your mail")
	})

	router.Use(func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return templates.NotFound().Render(c.Context(), c)
	})
}
