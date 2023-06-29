package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

/*
 * Handler to write a log into the database.
 * Expects Authorization header.
 *
 * Expected body:
 *   Message: string
 *   Severity: int
 *   Origin: string
 */
func (app *App) WriteLogHandler(c *fiber.Ctx) error {
	var log LogSchema

	// Limit body size to 1MB and disallow unknown JSON fields
	err := c.BodyParser(&log)
	if err != nil {
		// TODO: better error handling
		return fiber.ErrBadRequest
	}

	log.TimeWritten = time.Now().UTC()

	ctx, cancel := context.WithTimeout(
		c.Context(),
		time.Duration(15*time.Second),
	)
	defer cancel()

	// Insert schema into database.
	_, err = app.logger.Database.InsertOne(ctx, log)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

/*
 * Handler to read all logs. No body is necessary.
 * Expects Authorization header.
 *
 * Return schema:
 *   TimeWritten: datetime
 *   Message: string
 *   Category: string
 *   Severity: int
 */
func (app *App) ReadLogsHandler(c *fiber.Ctx) error {
	// bson.M{} applies no filter.
	filter := bson.M{}

	for k, v := range c.AllParams() {
		filter[k] = v[0]
	}

	logs, err := app.ReadLogs(c.Context(), filter)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Response().Header.Set("Content-Type", "application/json")

	// Write logs to client.
	err = c.JSON(logs)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

/*
 * Handler to read top priority categories.
 * Expects Authorization header.
 *
 * Return schema:
 *   TimeWritten: datetime
 *   Message: string
 *   Category: string
 *   Severity: int
 */
func (app *App) PrioritiesHandler(c *fiber.Ctx) error {
	categories, err := app.FindPriority(c.Context())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Response().Header.Set("Content-Type", "application/json")

	// Write logs to client.
	err = c.JSON(categories)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}
