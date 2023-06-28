package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CloseDatabase(c *mongo.Client, exitCode int) int {
	if err := c.Disconnect(context.Background()); err != nil {
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
	opts := options.
		Client().
		ApplyURI(*DatabaseUrl).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	// Ping database to confirm connection.
	err = client.
		Database(*DatabaseName).
		RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).
		Err()

	if err != nil {
		return nil, err
	}

	/*
	 * Handle graceful exit on SIGTERM or SIGINT.
	 * without this, the database connection won't close properly.
	 */
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func(db *mongo.Client) {
		signal := <-sigc
		fmt.Printf(
			"recieved signal %v. successfully closed connections to database",
			signal,
		)
		os.Exit(CloseDatabase(client, 0))
	}(client)

	return client, nil
}
