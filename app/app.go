package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/arael34/qlogger/types"
	"go.mongodb.org/mongo-driver/mongo"
)

// App ----------------------------------------------------------------------

type App struct {
	client *mongo.Client
	logger *types.QLogger
}

func (app *App) Run() error {
	// Set up routes
	http.HandleFunc("/api/write/", app.WriteLog)
	http.HandleFunc("/api/read/", app.ReadLogs)

	http.Handle("/", http.FileServer(http.Dir("./static/")))

	// For production
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("\nready to go.")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("error serving: %v", err)
		os.Exit(CloseDatabase(app.client, 1))
	}

	return nil
}

// App Builder ---------------------------------------------------------------

type AppBuilder struct {
	client *mongo.Client
	logger *types.QLogger
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func (ab *AppBuilder) WithClient(client *mongo.Client) *AppBuilder {
	ab.client = client
	return ab
}

func (ab *AppBuilder) WithLogger(logger *types.QLogger) *AppBuilder {
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
