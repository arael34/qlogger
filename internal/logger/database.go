package logger

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CloseDatabase(c *mongo.Client, exitCode int) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	if err := c.Disconnect(ctx); err != nil {
		panic(err)
	}

	// for use with os.Exit(CloseDatabase())
	return exitCode
}

func ConnectToDatabase(
	DatabaseUrl *string,
	DatabaseName *string,
) (*mongo.Client, error) {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(*DatabaseUrl).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*time.Second))
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Ping database to confirm connection.
	err = client.Database(*DatabaseName).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
