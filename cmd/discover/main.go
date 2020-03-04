package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/kevinwylder/sbvision"

	"github.com/kevinwylder/sbvision/cmd"
	"github.com/kevinwylder/sbvision/database"
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
		log.Fatal(err)
	}

	assets, tmpdir := cmd.GetLocalAssets()
	if tmpdir != "" {
		defer os.RemoveAll(tmpdir)
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
	err = info.GetThumbnail(video.Thumbnail(), assets)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success")

}
