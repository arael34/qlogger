package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/*
 * Handler to write a log into the database.
 * Expects Authorization header.
 *
 * Expected body:
 *   Message: string
 *   Severity: int
 *   Origin: string
 */
func (app *App) WriteLogHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Authorization") != *app.logger.AuthHeader {
		http.Error(w, "not authorized.", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var log LogSchema

	// Limit body size to 1MB and disallow unknown JSON fields
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1048576))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&log)
	if err != nil {
		// TODO: better error handling
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.TimeWritten = time.Now().UTC()

	ctx, cancel := context.WithTimeout(
		r.Context(),
		time.Duration(15*time.Second),
	)
	defer cancel()

	// Insert schema into database.
	_, err = app.logger.Database.InsertOne(ctx, log)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * Handler to read all logs. No body is necessary.
 * Expects Authorization header.
 *
 * Return schema:
 *   TimeWritten: datetime
 *   Message: string
 *   Category: string
 *   Severity: int
 */
func (app *App) ReadLogsHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Authorization") != *app.logger.AuthHeader {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// bson.M{} applies no filter.
	filter := bson.M{}

	for k, v := range r.URL.Query() {
		filter[k] = v[0]
	}

	logs, err := app.ReadLogs(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write logs to client.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * Handler to read top priority categories.
 * Expects Authorization header.
 *
 * Return schema:
 *   TimeWritten: datetime
 *   Message: string
 *   Category: string
 *   Severity: int
 */
func (app *App) PrioritiesHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Authorization") != *app.logger.AuthHeader {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := app.FindPriority(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write logs to client.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
