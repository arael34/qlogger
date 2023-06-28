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
