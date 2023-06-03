package app

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arael34/qlogger/types"
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
func (app *App) WriteLog(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Authorization") != *app.logger.AuthHeader {
		http.Error(w, "not authorized.", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var log types.LogSchema

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

	// If there is a current read connection open, write to it.
	if app.logger.Conn != nil {
		app.logger.Conn.WriteJSON(log)
	}
}

/*
 * Handler to read all logs. No body is necessary.
 * Expects Sec-WebSocket-Protocol header (for auth).
 *
 * Return schema:
 *   TimeWritten: datetime
 *   Message: string
 *   Origin: string
 *   Severity: int
 */
func (app *App) ReadLogs(w http.ResponseWriter, r *http.Request) {
	// Authorize user.
	if r.Header.Get("Sec-WebSocket-Protocol") != *app.logger.AuthHeader {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(
		r.Context(),
		time.Duration(15*time.Second),
	)
	defer cancel()

	// bson.D{} applies no filter.
	cursor, err := app.logger.Database.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse documents as schema.
	var logs []types.LogSchema
	err = cursor.All(ctx, &logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Connect to websocket for realtime updates.
	conn, err := app.logger.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	app.logger.Conn = conn

	for _, log := range logs {
		conn.WriteJSON(log)
	}

	app.logger.HandleSocket(conn)
}
