package main

import (
	"fmt"
	"log"
	"majin/config"
	"majin/internal/utils"
	"majin/middleware"
	"majin/router"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
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
			// Allow both http and https, and any subdomain
			fmt.Sprintf("http://*.%s", domain),
			fmt.Sprintf("https://*.%s", domain),

			// Also allow the root domain
			fmt.Sprintf("http://%s", domain),
			fmt.Sprintf("https://%s", domain),
		)
	}

	return strings.Join(origins, ",")
}

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if config.Debug {
		log.Printf("Running in debug mode on port %s", config.Port)
		log.Printf("Authorized domains: %v", config.AuthorizedDomains)
	}

	store := session.New(session.Config{
		KeyLookup:      "cookie:session_id",
		Expiration:     24 * time.Hour,
		CookieSecure:   !config.Debug,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	app := fiber.New()

	app.Use(recover.New(recover.Config{
		EnableStackTrace: config.Debug,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     buildCORSOrigins(config.AuthorizedDomains),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		ExposeHeaders:    "Content-Length, Content-Type",
		MaxAge:           0,
	}))
	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.GenerateToken,
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("session", store)
		return c.Next()
	})

	app.Static("/static", "./static")

	app.Use(middleware.HTMLMiddleware())
	app.Use(middleware.AuthMiddleware(store))

	router.Initialize(app)

	log.Fatal(app.Listen(":" + config.Port))
}
