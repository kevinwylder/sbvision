package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/kevinwylder/sbvision/video"
)

func (ctx *serverContext) handleGetVideo(w http.ResponseWriter, r *http.Request) {
	if r.Form.Get("id") == "" {
		http.Error(w, "?id= query param string required", 400)
		return
	}

	video, err := ctx.ddb.GetVideoByID(r.Form.Get("id"))
	if err != nil {
		http.Error(w, "Could not get video", 500)
		return
	}
	json.NewEncoder(w).Encode(video)
}

func (ctx *serverContext) handleVideoPage(w http.ResponseWriter, r *http.Request) {
	_, err := ctx.auth.User(r.Header.Get("Identity"))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	videos, err := ctx.ddb.GetAllVideos()
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

	requests := ctx.processes.GetUserRequests(user)

	url := r.Form.Get("url")
	if url != "" {
		requests.NewRequest(url, "", nil)
	} else {
		title := r.Form.Get("title")
		if title == "" {
			http.Error(w, "url is missing", 400)
			return
		}
		file, _, err := r.FormFile("video")
		if err != nil {
			http.Error(w, "video is missing", 400)
			return
		}
		requests.NewRequest("", title, file.(*os.File))
	}
}

func (ctx *serverContext) handleVideoStatus(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Form.Get("identity"))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	requests := ctx.processes.GetUserRequests(user)

	socket, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open socket", 500)
		return
	}

	var callbackID int64
	callbackID = requests.AddListener(func(status *video.Status) {
		err := socket.WriteJSON(status)
		if err != nil {
			requests.RemoveListener(callbackID)
			socket.Close()
		}

	})
}
