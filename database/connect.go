package database

import (
	"database/sql"
	"fmt"
)

// SBDatabase is a namespace of database queries
type SBDatabase struct {
	db               *sql.DB
	addImage         *sql.Stmt
	addSession       *sql.Stmt
	addVideo         *sql.Stmt
	addYoutubeRecord *sql.Stmt
	addFrame         *sql.Stmt
	getVideoPage     *sql.Stmt
	getVideoCount    *sql.Stmt
	getVideoByID     *sql.Stmt
	getYoutubeRecord *sql.Stmt
	getFrame         *sql.Stmt
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
	err = sb.prepareAddSession()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing AddSession: %s", err.Error())
	}
	err = sb.prepareAddImage()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing AddImage: %s", err.Error())
	}
	err = sb.prepareAddVideo()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing AddVideo: %s", err.Error())
	}
	err = sb.prepareAddYoutubeRecord()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing AddYoutubeRecord: %s", err.Error())
	}
	err = sb.prepareGetVideoCount()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing GetVideoCount: %s", err.Error())
	}
	err = sb.prepareGetVideos()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing GetVideos: %s", err.Error())
	}
	err = sb.prepareGetVideoByID()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing GetVideoByID: %s", err)
	}
	err = sb.prepareGetYoutubeRecord()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing GetYoutubeRecord: %s", err.Error())
	}
	err = sb.prepareGetFrame()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing GetFrame: %s", err)
	}
	err = sb.prepareAddFrame()
	if err != nil {
		return nil, fmt.Errorf("\n\tError prepareing AddFrame: %s", err)
	}
	return sb, nil
}
