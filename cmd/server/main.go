package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/cmd"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/frontend"
	"github.com/kevinwylder/sbvision/session"
	"github.com/kevinwylder/sbvision/youtube"
)

type serverContext struct {
	session  sbvision.SessionManager
	assets   sbvision.KeyValueStore
	youtube  sbvision.VideoHandler
	frontend http.Handler
	db       *database.SBDatabase
}

func main() {
	db, err := database.ConnectToDatabase(os.Getenv("DB_CREDS"))
	counter := 0
	for err != nil {
		log.Print(err)
		time.Sleep(time.Second)
		db, err = database.ConnectToDatabase(os.Getenv("DB_CREDS"))
		counter++
		if counter > 30 {
			log.Fatal("Could not connect to the database :(")
		}
	}

	session, err := session.NewRSASessionManager()
	if err != nil {
		log.Fatal(err)
	}

	server := &serverContext{
		db:      db,
		session: session,
	}

	if _, exists := os.LookupEnv("FRONTEND_DIR"); !exists {
		log.Fatal("Missing FRONTEND_DIR env variable")
	}
	server.frontend, err = frontend.ServeFrontend(os.Getenv("FRONTEND_DIR"))
	if err != nil {
		log.Fatal(err)
	}

	assets, cleanup := cmd.GetLocalAssets()
	if cleanup != "" {
		defer os.RemoveAll(cleanup)
	}
	server.assets = assets

	server.youtube = youtube.NewYoutubeHandler(db, server.assets)

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), server))
}
