package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleImage(w http.ResponseWriter, r *http.Request) {
	image := r.URL.Path[7:]
	switch r.Method {
	case http.MethodPost:
		session, err := ctx.session.ValidateSession(sbvision.SessionJWT(r.Header.Get("Session")))
		if err != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}

		err = ctx.assets.PutAsset(image, r.Body)
		if err != nil {
			fmt.Println("Error saving image", err)
			http.Error(w, "Could not save image", 500)
			return
		}

		err = ctx.db.AddImage(sbvision.Image(image), session)
		if err != nil {
			fmt.Println("Error storing image", err)
			http.Error(w, "Could not store image in db", 500)
			return
		}

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
