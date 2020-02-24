package main

import (
	"fmt"
	"io"
	"net/http"
)

func (ctx *serverContext) handleImage(w http.ResponseWriter, r *http.Request) {
	image := r.URL.Path[7:]
	switch r.Method {

	case http.MethodGet:
		data, err := ctx.assets.GetAsset(image)
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
