package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// App ----------------------------------------------------------------------

type App struct {
	client *mongo.Client
	logger *QLogger
}

func (app *App) Run() error {
	router := fiber.New(fiber.Config{BodyLimit: 1048576})

	// Set up routes
	router.Post("/api/write/", app.WriteLogHandler)
	router.Get("/api/read/", app.ReadLogsHandler)
	router.Get("/api/insights/priorities/", app.PrioritiesHandler)

	router.Static("/", "./static/")

	// For production
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("\nlistening...")

	err := router.Listen(":" + port)
	if err != nil {
		fmt.Printf("error serving: %v", err)
		os.Exit(CloseDatabase(app.client, 1))
	}

	return nil
}

// App Builder ---------------------------------------------------------------

type AppBuilder struct {
	client *mongo.Client
	logger *QLogger
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func (ab *AppBuilder) WithClient(client *mongo.Client) *AppBuilder {
	ab.client = client
	return ab
}

func (ab *AppBuilder) WithLogger(logger *QLogger) *AppBuilder {
	ab.logger = logger
	return ab
}

func (ab *AppBuilder) Build() (*App, error) {
	if ab.client == nil ||
		ab.logger == nil {
		return nil, errors.New("failed to build app")
	}

	return &App{
		client: ab.client,
		logger: ab.logger,
	}, nil
}
