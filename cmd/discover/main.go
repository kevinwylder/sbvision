package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kevinwylder/sbvision/cmd"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/youtube"

	"github.com/kevinwylder/sbvision"
)

func main() {

	flags := flag.NewFlagSet("discover", flag.ExitOnError)
	url := flags.String("url", "", "A URL to discover videos on")
	if err := flags.Parse(os.Args); err != nil {
		flags.PrintDefaults()
		log.Fatal(err)
	}

	db, err := database.ConnectToDatabase(os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err)
	}

	assets, tmpdir := cmd.GetLocalAssets()
	if tmpdir != "" {
		defer os.RemoveAll(tmpdir)
	}

	// Route to index a video
	var video sbvision.VideoDiscoverRequest
	video.Type = 1
	video.URL = *url

	youtube := youtube.NewYoutubeHandler(db, assets)

	// Only youtube is supported at this time, here is the "polymorphic dispatch"
	v, err := youtube.HandleDiscover(&video)
	if err != nil {
		log.Fatal("YoutubeHandler download error ", err.Error())
		return
	}

	fmt.Println("Successfully found", v)

}
