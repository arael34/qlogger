package main

import (
	"fmt"
	qlogger "internal/logger"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("starting...")
	fmt.Println()

	// Grab environment variables
	environment, envErr := ValidateEnvironment()
	if envErr != nil {
		fmt.Println(envErr)
		os.Exit(1)
	}

	fmt.Println("successfully loaded environment.")
	// Finish grabbing environment variables

	// Connect to database
	db, dbErr := qlogger.ConnectToDatabase(&environment.AuthHeader)
	if dbErr != nil {
		fmt.Println(envErr)
		os.Exit(1)
	}

	/*
	 * Handle graceful exit on SIGTERM or SIGINT.
	 * without this, the database connection won't close properly.
	 */
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func(db *qlogger.QLoggerDatabase) {
		signal := <-sigc
		db.Close()
		fmt.Printf(
			"recieved signal %v. successfully closed connections to database",
			signal,
		)
		os.Exit(0)
	}(db)

	// On normal exit, close connection to db.
	defer func(db *qlogger.QLoggerDatabase) {
		db.Close()
		fmt.Printf("successfully closed connection to database")
	}(db)

	if err := db.Handle.Ping(); err != nil {
		fmt.Printf("failed to ping: %v", err)
		os.Exit(db.Close())
	}

	fmt.Println("successfully pinged database.")
	// Finish connecting to database

	fmt.Println("\nready to go.")
	fmt.Println()

	// Set up routes
	qlog := qlogger.NewQLogger(&environment.AuthHeader, db)

	http.HandleFunc("/api/write/", qlog.WriteLog)
	http.HandleFunc("/api/read/", qlog.ReadLog)

	http.Handle("/", http.FileServer(http.Dir("../static")))
	// Finish setting up routes

	serveErr := http.ListenAndServe(":3000", nil)
	if serveErr != nil {
		fmt.Printf("error serving: %v", serveErr)
		os.Exit(db.Close())
	}
}
