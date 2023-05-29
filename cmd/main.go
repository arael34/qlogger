package main

import (
	"fmt"
	qlogger "internal/logger"
	"net/http"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	fmt.Println("starting...")
	fmt.Println()

	// Grab environment variables
	environment, err := qlogger.ValidateEnvironment()
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("successfully loaded environment.")
	// Finish grabbing environment variables

	// Connect to database
	client, err := qlogger.ConnectToDatabase(
		&environment.DatabaseUrl,
		&environment.DatabaseName,
	)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
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
		os.Exit(qlogger.CloseDatabase(client, 0))
	}(client)

	fmt.Println("successfully pinged database.")
	// Finish connecting to database

	// Set up routes
	qlog := qlogger.NewQLogger(
		&environment.AuthHeader,
		client.Database("qlogger").Collection("logs"),
	)

	http.HandleFunc("/api/write/", qlog.WriteLog)
	http.HandleFunc("/api/read/", qlog.ReadLogs)

	http.Handle("/", http.FileServer(http.Dir("./static/")))
	// Finish setting up routes

	// For production
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("\nready to go.")
	fmt.Println()

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("error serving: %v", err)
		os.Exit(qlogger.CloseDatabase(client, 1))
	}
}
