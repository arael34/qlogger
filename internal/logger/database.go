package logger

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type QLoggerDatabase struct {
	Handle   *sql.DB
	WriteLog *sql.Stmt
	ReadLog  *sql.Stmt
}

func (q *QLoggerDatabase) Close() int {
	q.Handle.Close()
	q.WriteLog.Close()
	q.ReadLog.Close()

	// for use with os.Exit(db.Close())
	return 1
}

func ConnectToDatabase(DatabaseUrl *string) (*QLoggerDatabase, error) {
	Handle, handleErr := sql.Open("mysql", *DatabaseUrl)
	if handleErr != nil {
		return nil, errors.New("error connecting to database.")
	}

	// Important settings
	Handle.SetConnMaxLifetime(time.Minute * 3)
	Handle.SetMaxOpenConns(10)
	Handle.SetMaxIdleConns(10)

	// Prepare statements
	WriteLog, wErr := Handle.Prepare("INSERT INTO logs VALUES ( ?, ?, ? )")
	if wErr != nil {
		return nil, errors.New("error preparing write handle.")
	}

	ReadLog, rErr := Handle.Prepare("SELECT message FROM logs WHERE message = ?")
	if rErr != nil {
		return nil, errors.New("error preparing read handle.")
	}

	return &QLoggerDatabase{Handle, WriteLog, ReadLog}, nil
}
