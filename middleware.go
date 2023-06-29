package main

import "github.com/gofiber/fiber/v2"

type QMiddlewareConfig struct {
	Filter      func(c *fiber.Ctx) bool
	MaxBodySize int64
	AuthHeader  string
	Origins     map[string]bool
}

func QMiddleware(config QMiddlewareConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Filter != nil && !config.Filter(c) {
			return c.Next()
		}

		// Check authorization header
		if config.AuthHeader != "" {
			auth := c.Get("Authorization")
			if auth == config.AuthHeader {
				return c.Status(401).SendString("Unauthorized")
			}
		}

		// Check origin
		if len(config.Origins) > 0 {
			origin := c.Get("Origin")
			if !config.Origins[origin] {
				return c.Status(403).SendString("Forbidden")
			}
		}

		return c.Next()
	}
}
