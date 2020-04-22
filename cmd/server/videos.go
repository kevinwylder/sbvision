package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kevinwylder/sbvision"

	"github.com/kevinwylder/sbvision/media/sources"
)

func (ctx *serverContext) handleGetVideo(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, []string{"id"})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	video, err := ctx.db.GetVideoByID(ids[0])
	if err != nil {
		http.Error(w, "Could not get video", 500)
		return
	}
	json.NewEncoder(w).Encode(video)
}

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

func (ctx *serverContext) handleVideoDiscovery(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Header.Get("Identity"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", 401)
		return
	}

	err = r.ParseMultipartForm(2048)
	if err != nil {
		http.Error(w, "could not parse multipart form", 400)
		return
	}

	url := r.Form.Get("url")
	if url != "" {
		ticket, err = ctx.discoveryQueue.Enqueue(user, func() (sbvision.VideoSource, error) {
			return sources.FindVideoSource(url)
		})
	} else {
		title := r.Form.Get("title")
		if title == "" {
			http.Error(w, "either (url) or (title, video) required in form data", 400)
			return
		}
		file, _, err := r.FormFile("video")
		if err != nil {
			http.Error(w, "video missing from form", 400)
			return
		}
		ticket, err = ctx.discoveryQueue.Enqueue(user, func() (sbvision.VideoSource, error) {
			return sources.VideoFileSource(file, title, func() {
				r.MultipartForm.RemoveAll()
			})
		})

	}
}

func (ctx *serverContext) handleVideoStatus(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Form.Get("identity"))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	socket, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open socket", 500)
		return
	}

	request, exists := ctx.discoveryQueue.Find(user)
	if !exists {
		socket.Close()
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
