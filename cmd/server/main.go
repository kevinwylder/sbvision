package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/auth"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/media"
	"github.com/kevinwylder/sbvision/media/video"

	"github.com/gorilla/websocket"
)

type serverContext struct {
	assets         *media.AssetDirectory
	upgrader       websocket.Upgrader
	auth           *auth.JWTVerifier
	discoveryQueue *video.ProcessQueue
	videoCache     map[int64]*sbvision.Video
	db             *database.SBDatabase
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

	server := &serverContext{
		db: db,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 20 * 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		auth:       auth.NewJWTVerifier(db, "https://cognito-idp.us-west-2.amazonaws.com/us-west-2_dHWlJDm4T/.well-known/jwks.json"),
		videoCache: make(map[int64]*sbvision.Video),
	}

	if server.assets, err = media.NewAssetDirectory(os.Getenv("ASSET_DIR")); err != nil {
		log.Fatal(err)
	}
	server.discoveryQueue = video.NewProcessQueue(server.assets, db)

	fmt.Println("Starting server")
	err = http.ListenAndServe(":"+os.Getenv("PORT"), server)
	fmt.Println(err)
}
