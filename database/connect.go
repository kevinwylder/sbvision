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
	getVideoPage     *sql.Stmt
	getVideoCount    *sql.Stmt
	getYoutubeRecord *sql.Stmt
}

// ConnectToDatabase uses the DB_CREDS environment variable to connect to the database
func ConnectToDatabase(creds string) (*SBDatabase, error) {
	db, err := sql.Open("mysql", creds)
	if err != nil {
		return nil, fmt.Errorf("\n\tError connecting to the database: %s", err.Error())
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
	err = sb.prepareGetYoutubeRecord()
	if err != nil {
		return nil, fmt.Errorf("\n\tError preparing GetYoutubeRecord: %s", err.Error())
	}
	return sb, nil
}
