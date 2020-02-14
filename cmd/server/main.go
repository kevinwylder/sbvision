package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/assets/amazon"
	"github.com/kevinwylder/sbvision/assets/filesystem"
	"github.com/kevinwylder/sbvision/cropper"
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
	cropper  *cropper.PngCropper
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

	assetDir, exists := os.LookupEnv("ASSET_DIR")
	if !exists {
		assetDir, err = ioutil.TempDir("", "")
		if err != nil {
			log.Fatal("Could not create tmp dir for image storage", err)
		}
	}
	cache, err := filesystem.NewAssetDirectory(assetDir)
	if err != nil {
		log.Fatal(err)
	}
	if bucket, exists := os.LookupEnv("S3_BUCKET"); exists {
		server.assets, err = amazon.NewS3BucketManager(bucket, cache)
		if err != nil {
			log.Fatal(err)
		}
	}
	if server.assets == nil {
		server.assets = cache
	}

	server.cropper = cropper.NewPngCropper(server.assets)

	server.youtube = youtube.NewYoutubeHandler(db, server.assets)

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), server))
}
