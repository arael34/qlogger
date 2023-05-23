package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("starting...")

	// Grab environment variables
	environment, envErr := qlogger.ValidateEnvironment()
	if envErr != nil {
		fmt.Println(envErr)
		os.Exit(1)
	}

	fmt.Println("successfully loaded environment.")
	// Finish grabbing environment variables

	// Connect to database
	db, dbErr := sql.Open("mysql", environment.DatabaseUrl)
	if dbErr != nil {
		fmt.Println(dbErr)
		os.Exit(1)
	}

	// Important settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	/*
	 * Handle graceful exit on SIGTERM or SIGINT.
	 * without this, the database connection won't close properly.
	 */
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func(db *sql.DB) {
		signal := <-sigc
		db.Close()
		fmt.Printf(
			"recieved signal %v. successfully closed connection to database",
			signal,
		)
		os.Exit(0)
	}(db)

	// On normal exit, close connection to db.
	defer func(db *sql.DB) {
		db.Close()
		fmt.Printf("successfully closed connection to database")
	}(db)

	if err := db.Ping(); err != nil {
		fmt.Printf("failed to ping: %v", err)
		os.Exit(1)
	}

	fmt.Println("successfully pinged database.")
	// Finish connecting to database

	// Set up routes
	qlog := qlogger.NewQLogger(&environment.DatabaseUrl)

	http.HandleFunc("/api/write/", qlog.WriteLog)
	http.HandleFunc("/api/read/", qlog.ReadLog)

	http.Handle("/", http.FileServer(http.Dir("../static")))
	// Finish setting up routes

	log.Fatal(http.ListenAndServe(":3000", nil))
}
