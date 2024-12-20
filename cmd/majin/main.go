package main

import (
	"fmt"
	"log"
	"majin/config"
	"majin/middleware"
	"majin/router"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func buildCORSOrigins(domains []string) string {
	origins := []string{
		"http://localhost",
		"https://localhost",
	}

	for _, domain := range domains {
		origins = append(origins,
			fmt.Sprintf("http://*.%s", domain),
			fmt.Sprintf("https://*.%s", domain),
			fmt.Sprintf("http://%s", domain),
			fmt.Sprintf("https://%s", domain),
		)
	}

	return strings.Join(origins, ",")
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	store := session.New(session.Config{
		KeyLookup:      "cookie:session_id",
		Expiration:     24 * time.Hour,
		CookieSecure:   !cfg.Debug,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("session", store)
		return c.Next()
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.Debug,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     buildCORSOrigins(cfg.AuthorizedDomains),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Csrf-Token",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		ExposeHeaders:    "Content-Length, Content-Type, X-Csrf-Token",
	}))

	// Conditional logging based on debug mode
	if cfg.Debug {
		app.Use(logger.New())
	}

	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy:     "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: blob:; connect-src 'self'; frame-ancestors 'none';",
		XFrameOptions:             "SAMEORIGIN",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		CrossOriginEmbedderPolicy: "credentialless",
		CrossOriginResourcePolicy: "same-site",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	})

	app.Static("/static", "./static")

	app.Use(middleware.CSRFMiddleware(!cfg.Debug))
	app.Use(middleware.HTMLMiddleware())
	app.Use(middleware.AuthMiddleware(store))
	app.Use(middleware.TemplateContextMiddleware())

	router.Initialize(app)

	log.Fatal(app.Listen(":" + cfg.Port))
}
