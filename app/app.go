package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arael34/qlogger/types"
	"go.mongodb.org/mongo-driver/mongo"
)

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
