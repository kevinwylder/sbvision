package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/kevinwylder/sbvision"
)

type serverContext struct {
	db                *sql.DB
	sessionManager    sbvision.SessionManager
	sessionStorage    sbvision.SessionStorage
	youtubeDownloader sbvision.VideoDownloader
	videoLister       sbvision.VideoList
}

func main() {
	db, err := sql.Open("mysql", os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err)
	}

	server := &serverContext{db}

	http.HandleFunc("/session", server.session)
	http.HandleFunc("/videos", server.videos)
}
