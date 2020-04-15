package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/media/video"

	"github.com/kevinwylder/sbvision/media/sources"
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

func (ctx *serverContext) handleVideoStreamInit(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.FormValue("user"))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	id, err := getIDs(r, []string{"id"})
	if err != nil {
		http.Error(w, "Missing video id", 400)
		return
	}

	video, uploader, err := ctx.db.GetVideoByID(id[0])
	if err != nil {
		http.Error(w, "Cannot get video", 500)
		return
	}

	if uploader != user.ID {
		http.Error(w, "Not your video", 401)
		return
	}

	random := make([]byte, 20)
	rand.Read(random)
	key := base64.URLEncoding.EncodeToString(random)
	ctx.videoCache[key] = video

	w.Header().Set("Location", "/video/stream?key="+key)
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(307)
}

func (ctx *serverContext) handleVideoStream(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "Missing key", 400)
		return
	}
	video, exists := ctx.videoCache[key]
	if !exists {
		http.Error(w, "Not Found", 404)
		return
	}
	http.ServeFile(w, r, ctx.assets.VideoPath(video))
}

func (ctx *serverContext) handleVideoDiscovery(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Header.Get("Identity"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", 401)
		return
	}

	err = r.ParseMultipartForm(1024 * 5)
	if err != nil {
		http.Error(w, "could not parse multipart form", 400)
		return
	}

	url := r.Form.Get("url")
	var ticket *video.ProcessRequest
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
			return sources.VideoFileSource(file, title)
		})

	}
	if err != nil {
		http.Error(w, err.Error(), 400)
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
