package database

import (
	// this package only supports mysql
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"reflect"
	"runtime"
)

// SBDatabase is a namespace of database queries
type SBDatabase struct {
	db *sql.DB

	// in session.go
	addSession *sql.Stmt

	// in contribute.go
	addFrame    *sql.Stmt
	addBounds   *sql.Stmt
	addRotation *sql.Stmt

	// in dataset.go
	dataWhereVideo      *sql.Stmt
	dataWhereFrame      *sql.Stmt
	dataWhereNoRotation *sql.Stmt

	// in video.go
	addVideo      *sql.Stmt
	updateVideo   *sql.Stmt
	getVideoByID  *sql.Stmt
	getVideoPage  *sql.Stmt
	getVideoCount *sql.Stmt
}

// ConnectToDatabase waits for a sql connection then prepares queries for runtime
func ConnectToDatabase(creds string) (*SBDatabase, error) {
	db, err := sql.Open("mysql", creds)
	if err != nil {
		return nil, fmt.Errorf("\n\tError connecting to the database: %s", err.Error())
	}
	var counter int
	for err := db.Ping(); err != nil; counter++ {
		if counter > 30 {
			return nil, fmt.Errorf("Timed out pinging the database: %s", err.Error())
		}
	}

	sb := &SBDatabase{db: db}

	var prepFunctions = [](func() error){
		// in session.go
		sb.prepareAddSession,

		// in contribute.go
		sb.prepareAddFrame,
		sb.prepareAddBounds,
		sb.prepareAddRotation,

		// in dataset.go
		sb.prepareDataWhereVideo,
		sb.prepareDataWhereFrame,
		sb.prepareDataWhereNoRotation,

		// in video.go
		sb.prepareAddVideo,
		sb.prepareGetVideoByID,
		sb.prepareUpdateVideo,
		sb.prepareGetVideos,
		sb.prepareGetVideoCount,
	}

	for _, f := range prepFunctions {
		err = f()
		if err != nil {
			return nil, fmt.Errorf("\n\tError preparing query %s: %s", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), err.Error())
		}
	}
	return sb, nil
}

// this interface is used to scan more generic built queries
type scannable interface {
	Scan(to ...interface{}) error
}
