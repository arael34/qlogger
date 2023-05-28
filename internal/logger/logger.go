package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
 * Simple alias for readability.
 * 0 - DEBUG
 * 1 - WARN
 * 2 - ERROR
 */
type Level int

const (
	DEBUG Level = iota
	WARN
	ERROR
)

type QLogger struct {
	authHeader *string
	database   *mongo.Collection
	upgrader   *websocket.Upgrader
	conn       *websocket.Conn
}

func NewQLogger(authHeader *string, database *mongo.Collection) *QLogger {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO fix CheckOrigin
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return &QLogger{authHeader, database, upgrader, nil}
}

/*
 * Schema for a single log entry.
 */
type LogSchema struct {
	TimeWritten time.Time `bson:"time"`
	Message     string    `bson:"message"`
	Origin      string    `bson:"origin"`
	Severity    Level     `bson:"severity"`
}

func (logger *QLogger) HandleSocket(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		if err = conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}
