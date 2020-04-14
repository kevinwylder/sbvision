package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (ctx *serverContext) handleVideoPage(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Header.Get("Identity"))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	videos, err := ctx.db.GetVideos(user)
	if err != nil {
		http.Error(w, "An error occured", 500)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(videos)
}

func (ctx *serverContext) handleVideoThumbnail(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, []string{"id"})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	data, err := ctx.assets.GetThumbnail(ids[0])
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

func (ctx *serverContext) handleVideoStream(w http.ResponseWriter, r *http.Request) {

}

func (ctx *serverContext) handleVideoDiscovery(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Header.Get("Identity"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", 401)
		return
	}

	var request struct {
		URL string `json:"url"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	ticket, err := ctx.discoveryQueue.Enqueue(user, request.URL)
	if err != nil {
		http.Error(w, "Queue is full, come back later", 503)
		return
	}
	json.NewEncoder(w).Encode(ticket)
}

func (ctx *serverContext) handleVideoStatus(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Form.Get("identity"))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	request, exists := ctx.discoveryQueue.Find(user)
	if !exists {
		http.Error(w, "Not Found", 404)
		return
	}

	socket, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open socket", 500)
		return
	}

	events, done := request.Subscribe()
	defer done()

	ticker := time.NewTicker(time.Second * 5)
	for {
		var err error
		select {
		case <-events:
			err = socket.WriteJSON(request)
		case <-ticker.C:
			err = socket.WriteJSON(request)
		}
		if err != nil {
			socket.Close()
		}
	}
}
