package lib

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type _QLogger struct {
	AuthHeader  *string
	Database    *mongo.Collection
	Upgrader    *websocket.Upgrader
	ConnSync    sync.Mutex
	Connections []*websocket.Conn
}

func _NewQLogger(authHeader *string, database *mongo.Collection) *_QLogger {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO fix CheckOrigin
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return &_QLogger{authHeader, database, upgrader, sync.Mutex{}, nil}
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

func (logger *_QLogger) CleanUp(conn *websocket.Conn) {
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
