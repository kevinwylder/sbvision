package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/cmd"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/sbvideo"
)

func main() {
	fmt.Println("Scraping reddit")

	db, err := database.ConnectToDatabase(os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err)
	}

	assets, tmpdir := cmd.GetLocalAssets()
	if tmpdir != "" {
		defer os.RemoveAll(tmpdir)
	}

	posts, err := sbvideo.GetRedditSkateboardingPosts()
	if err != nil {
		log.Fatal(err)
	}

	var video sbvision.Video
	for _, url := range posts {
		post, err := sbvideo.GetRedditPost(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		video.OriginURL = url
		post.Update(&video)
		if video.URL == "" {
			continue
		}
		fmt.Println("Found", video.Title)
		err = db.AddVideo(&video)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = post.GetThumbnail(video.Thumbnail(), assets)
		if err != nil {
			fmt.Println(err)
		}
	}

}
