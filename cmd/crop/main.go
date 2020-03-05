package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/kevinwylder/sbvision/sbimage"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/cmd"

	"github.com/kevinwylder/sbvision/database"
)

func main() {
	db, err := database.ConnectToDatabase(os.Getenv("DB_CREDS"))
	if err != nil {
		log.Fatal(err)
	}

	assets, tmp := cmd.GetLocalAssets()
	if tmp != "" {
		defer os.RemoveAll(tmp)
	}

	var offset int64
	for {
		page, err := db.DataWhereHasBound(offset)
		if err != nil {
			log.Fatal(err)
		}
		offset = page.NextOffset
		for i := range page.Frames {
			ensureCropped(assets, &page.Frames[i])
		}
		if !page.IsTruncated {
			return
		}
	}
}

func ensureCropped(assets sbvision.KeyValueStore, frame *sbvision.Frame) {
	var cropper *sbimage.PngCropper
	for i := range frame.Bounds {
		// first try to get the bounds image directly
		image, err := assets.GetAsset(frame.Bounds[i].Key())
		if err != nil {
			if cropper == nil {
				image, err = assets.GetAsset(frame.Key())
				if err != nil {
					fmt.Println("Missing frame", frame.ID, err)
					break
				}
				cropper, err = sbimage.Crop(image)
				image.Close()
				if err != nil {
					fmt.Println("Corrupted frame", frame.ID, err)
					break
				}
			}
			var buffer bytes.Buffer
			err = cropper.GetCroppedPng(&frame.Bounds[i], &buffer)
			if err != nil {
				fmt.Println("Error cropping bound", frame.Bounds[i].ID, err)
			}
			err = assets.PutAsset(frame.Bounds[i].Key(), &buffer)
			if err != nil {
				fmt.Println("Error storing cropped bound", frame.Bounds[i].ID)
			}
			fmt.Println("Cropped ", frame.Bounds[i].ID)
		} else {
			image.Close()
		}
	}
}
