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
	db               *sql.DB
	addImage         *sql.Stmt
	addSession       *sql.Stmt
	addVideo         *sql.Stmt
	addYoutubeRecord *sql.Stmt
	addFrame         *sql.Stmt
	addBounds        *sql.Stmt
	addRotation      *sql.Stmt

	getVideoPage     *sql.Stmt
	getVideoCount    *sql.Stmt
	getVideoByID     *sql.Stmt
	getYoutubeRecord *sql.Stmt
	getFrame         *sql.Stmt

	dataCounts         *sql.Stmt
	dataVideoFrames    *sql.Stmt
	dataAllFrames      *sql.Stmt
	dataRotationFrames *sql.Stmt
	dataByBoundID      *sql.Stmt

	setFrameHash *sql.Stmt
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
		sb.prepareAddBounds,
		sb.prepareAddFrame,
		sb.prepareAddImage,
		sb.prepareAddSession,
		sb.prepareAddVideo,
		sb.prepareAddYoutubeRecord,
		sb.prepareAddRotation,
		sb.prepareGetFrame,
		sb.prepareGetVideoByID,
		sb.prepareGetVideoCount,
		sb.prepareGetVideos,
		sb.prepareGetYoutubeRecord,
		sb.prepareDataVideoFrames,
		sb.prepareDataCounts,
		sb.prepareDataRotationFrames,
		sb.prepareDataByBoundID,
		sb.prepareSetFrameHash,
		sb.prepareDataAllFrames,
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
