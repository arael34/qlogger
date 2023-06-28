package types

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
 * Simple alias for readability.
 * 0 - INFO
 * 1 - DEBUG
 * 2 - WARN
 * 3 - ERROR
 */
type Level int

const (
	INFO Level = iota
	DEBUG
	WARN
	ERROR
)

type QLogger struct {
	AuthHeader  *string
	Database    *mongo.Collection
	Upgrader    *websocket.Upgrader
	ConnSync    sync.Mutex
	Connections []*websocket.Conn
}

func NewQLogger(authHeader *string, database *mongo.Collection) *QLogger {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO fix CheckOrigin
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return &QLogger{authHeader, database, upgrader, sync.Mutex{}, nil}
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

func (logger *QLogger) CleanUp(conn *websocket.Conn) {
	logger.ConnSync.Lock()
	defer logger.ConnSync.Unlock()

	for i, c := range logger.Connections {
		if c == conn {
			logger.Connections = append(
				logger.Connections[:i],
				logger.Connections[i+1:]...,
			)
		}
	}
}
