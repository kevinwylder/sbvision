package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/kevinwylder/sbvision"

	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/media"
	"github.com/kevinwylder/sbvision/sbvideo"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("USAGE: discover [url]")
	}
	discover, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	video := &sbvision.Video{
		OriginURL: os.Args[1],
	}

	db, err := database.ConnectToDatabase(os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err, "check the DB_CREDS environment variable")
	}

	assets, err := media.NewAssetDirectory(os.Getenv("ASSET_DIR"))
	if err != nil {
		log.Fatal(err, "check the ASSET_DIR environment variable")
	}

	fmt.Println("Getting info")
	var info sbvision.VideoSource

	switch discover.Host {
	case "www.youtube.com":
		info, err = sbvideo.GetYoutubeDl(video.OriginURL)
	case "www.reddit.com":
		info, err = sbvideo.GetRedditPost(video.OriginURL)
	default:
		log.Fatal("Cannot discover videos at", discover.Host)
	}

	if err != nil {
		log.Fatal(err)
	}

	info.Update(video)

	fmt.Println("Adding video to db")
	err = db.AddVideo(video)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Getting thumbnail")
	data, err := info.GetThumbnail()
	if err != nil {
		log.Fatal("Error getting thumbnail", err)
	}

	fmt.Println("Storing thumbnail")
	err = assets.PutThumbnail(video.ID, data)
	if err != nil {
		log.Fatal("Error storing thumbnail", err)
	}

	fmt.Println("success")

}
