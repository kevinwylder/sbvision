package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/images"
	"github.com/kevinwylder/sbvision/session"
	"github.com/kevinwylder/sbvision/youtube"
)

type serverContext struct {
	session  sbvision.SessionManager
	images   ImageManager
	youtube  sbvision.VideoHandler
	frontend http.Handler
	db       *database.SBDatabase
}

func main() {
	db, err := database.ConnectToDatabase(os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err)
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
	server.frontend = http.FileServer(FileSystem{http.Dir(os.Getenv("FRONTEND_DIR"))})

	if imageDir, exists := os.LookupEnv("IMAGE_DIR"); exists {
		server.images, err = images.NewImageDirectory(imageDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	if bucket, exists := os.LookupEnv("S3_BUCKET"); exists {
		server.images, err = images.NewImageBucket(bucket)
		if err != nil {
			log.Fatal(err)
		}
	}

	server.youtube = youtube.NewYoutubeHandler(db, server.images)

	http.ListenAndServe(os.Getenv("PORT"), server)
}

func (ctx *serverContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/videos") {
		ctx.videos(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/video") {
		ctx.video(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/image") {
		ctx.image(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/session") {
		ctx.getSession(w, r)
		return
	}
	ctx.frontend.ServeHTTP(w, r)
}
