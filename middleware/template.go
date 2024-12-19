package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const fiberCtxKey contextKey = "fiberCtx"

func TemplateContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.Context(), fiberCtxKey, c)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
