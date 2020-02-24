package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/kevinwylder/sbvision"

	"github.com/kevinwylder/sbvision/assets/amazon"
	"github.com/kevinwylder/sbvision/assets/filesystem"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/sbimage"
)

func main() {
	db, err := database.ConnectToDatabase(os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal("Could not connect to the db: " + err.Error())
	}
	images, err := db.DataAllFrames()
	if err != nil {
		log.Fatal("Could not get all images: " + err.Error())
	}
	videos, err := db.GetVideos(0, 100)
	if err != nil {
		log.Fatal("Could not get videos: " + err.Error())
	}

	assetDir, exists := os.LookupEnv("ASSET_DIR")
	if !exists {
		assetDir, err = ioutil.TempDir("", "")
		if err != nil {
			log.Fatal("Could not create tmp dir for image storage: " + err.Error())
		}
	}
	cache, err := filesystem.NewAssetDirectory(assetDir)
	if err != nil {
		log.Fatal("Could not create tmp dir: " + err.Error())
	}
	var assets sbvision.KeyValueStore
	if bucket, exists := os.LookupEnv("S3_BUCKET"); exists {
		assets, err = amazon.NewS3BucketManager(bucket, cache)
		if err != nil {
			log.Fatal("Could not open S3 bucket: " + err.Error())
		}
	}
	if assets == nil {
		assets = cache
	}

	var wg sync.WaitGroup
	for i := range images {
		wg.Add(1)
		go func(f *sbvision.Frame) {
			defer wg.Done()
			originalKey := fmt.Sprintf("frame/%d-%d.png", f.VideoID, f.Time)
			newKey := fmt.Sprintf("frame/%d.png", f.ID)
			data, err := assets.GetAsset(originalKey)
			if err != nil {
				fmt.Println("Could not get asset: " + err.Error())
				return
			}
			hash, err := sbimage.HashImage(data)
			data.Close()
			if err != nil {
				fmt.Println("Could not hash image " + originalKey + ": " + err.Error())
				return
			}
			err = db.SetFrameHash(hash, f.ID)
			if err != nil {
				fmt.Println(hash)
				fmt.Println("Could not set frame hash for " + originalKey + ": " + err.Error())
				return
			}

			data, err = assets.GetAsset(originalKey)
			if err != nil {
				fmt.Println("Error reopening", originalKey)
				return
			}
			defer data.Close()
			err = assets.PutAsset(newKey, data)
			if err != nil {
				fmt.Println("Error putting new key back: ", err.Error())
			}
		}(&images[i])
	}

	for i := range videos {
		wg.Add(1)
		go func(v *sbvision.Video) {
			defer wg.Done()
			data, err := assets.GetAsset(string(v.Thumbnail))
			if err != nil {
				fmt.Println("Error getting video asset", v.Thumbnail)
				return
			}
			defer data.Close()
			err = assets.PutAsset(fmt.Sprintf("thumbnail/%d.jpg", v.ID), data)
			if err != nil {
				fmt.Println("Error putting asset", err.Error())
			}
		}(&videos[i])
	}
	wg.Wait()
}
