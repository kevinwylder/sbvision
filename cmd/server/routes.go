package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/s3"
	"github.com/kevinwylder/sbvision/session"
	"github.com/kevinwylder/sbvision/youtube"
)

type serverContext struct {
	session sbvision.SessionManager
	images  *s3.ImageBucket
	youtube sbvision.VideoHandler
	db      *database.SBVisionDatabase
}

func main() {

	db, err := sql.Open("mysql", os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err)
	}

	session, err := session.NewRSASessionManager(db)
	if err != nil {
		log.Fatal(err)
	}

	bucket, err := s3.NewImageBucket()
	if err != nil {
		log.Fatal(err)
	}

	youtube := youtube.NewYoutubeHandler(db, bucket)

	server := &serverContext{
		db:      &database.SBVisionDatabase{db},
		images:  bucket,
		session: session,
		youtube: youtube,
	}

	http.HandleFunc("/session", server.getSession)
	http.HandleFunc("/videos", server.videos)
	http.HandleFunc("/video", server.video)
	http.ListenAndServe(":1080", nil)
}
