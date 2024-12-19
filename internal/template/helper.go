// internal/template/helper.go
package template

import (
	"context"
	"majin/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type TemplateData struct {
	CSRFToken string
}

func GetData(ctx context.Context) TemplateData {
	if fctx, ok := ctx.Value("fiberCtxKey").(*fiber.Ctx); ok {
		// Prioritize locals, then cookie
		csrfToken := fctx.Locals("_csrf")
		if csrfToken != nil {
			return TemplateData{
				CSRFToken: csrfToken.(string),
			}
		}

		cookieToken := fctx.Cookies("_csrf")
		if cookieToken != "" {
			return TemplateData{
				CSRFToken: cookieToken,
			}
		}
	}

	// Generate a new token if none exists
	return TemplateData{
		CSRFToken: utils.GenerateToken(),
	}
}
