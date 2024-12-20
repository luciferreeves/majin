package router

import (
	"majin/controllers"
	"majin/internal/auth"
	"majin/templates"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/mail")
	})

	router.Get("/login", controllers.LoginView)

	authRouter := router.Group("/auth")
	authRouter.Post("/login", controllers.Login)
	authRouter.Post("/logout", controllers.Logout)

	mail := router.Group("/mail", auth.RequireAuth)
	mail.Get("/", func(c *fiber.Ctx) error {
		return templates.Mail().Render(c.Context(), c.Response().BodyWriter())
	})

	router.Use(func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return templates.NotFound().Render(c.Context(), c)
	})
}
