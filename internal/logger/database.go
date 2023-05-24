package logger

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CloseDatabase(c *mongo.Client, exitCode int) int {
	if err := c.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	// for use with os.Exit(CloseDatabase())
	return exitCode
}

func ConnectToDatabase(DatabaseUrl *string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(*DatabaseUrl).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Ping database to confirm connection.
	// The database name should prob be in the env
	err = client.Database("qlogger").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
