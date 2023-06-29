package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (app *App) ReadLogs(
	ctx context.Context,
	filter bson.M,
) ([]LogSchema, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(15*time.Second),
	)
	defer cancel()

	cursor, err := app.logger.Database.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Parse documents as schema.
	var logs []LogSchema
	err = cursor.All(ctx, &logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (app *App) WriteLog(
	ctx context.Context,
	log LogSchema,
) error {
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(15*time.Second),
	)
	defer cancel()

	// Insert schema into database.
	_, err := app.logger.Database.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) DeleteLogs(
	ctx context.Context,
	filter bson.M,
) error {
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(15*time.Second),
	)
	defer cancel()

	_, err := app.logger.Database.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) FindPriority(ctx context.Context) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(15*time.Second),
	)
	defer cancel()

	cursor, err := app.logger.Database.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	categories := make(map[string]int)
	total := 0

	// Sift through logs.
	for cursor.Next(ctx) {
		var log LogSchema
		err := cursor.Decode(&log)
		if err != nil {
			return nil, err
		}

		if log.Severity == ERROR {
			total++
			if _, ok := categories[log.Category]; !ok {
				categories[log.Category] = 0
			}
		}
	}

	// Calculate percentages.
	for category := range categories {
		categories[category] = categories[category] * 100 / total
	}

	return categories, nil
}
