package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

// ImageManager is an uploader and downloader
type ImageManager interface {
	sbvision.ImageDownloader
	sbvision.ImageUploader
}

func (ctx *serverContext) image(w http.ResponseWriter, r *http.Request) {
	image := sbvision.Image(r.URL.Path[7:])
	switch r.Method {
	case http.MethodPost:
		err := ctx.images.PutImage(r.Body, image)
		if err != nil {
			fmt.Println("Error saving image", err)
			http.Error(w, "Could not save image", 500)
			return
		}

	case http.MethodGet:
		data, err := ctx.images.GetImage(image)
		if err != nil {
			fmt.Println("Could not get image", err)
			http.Error(w, "could not get image", 404)
			return
		}
		defer data.Close()

		_, err = io.Copy(w, data)
		if err != nil {
			fmt.Println("Error writing image response", err)
		}

	}
}
