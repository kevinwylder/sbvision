package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleVideoDiscovery(w http.ResponseWriter, r *http.Request) {
	// Route to index a video
	var video sbvision.VideoDiscoverRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&video)
	if err != nil {
		fmt.Println("VideoDownloadRequest decode error:", err.Error())
		http.Error(w, "invalid video format", 400)
		return
	}

	video.Session, err = ctx.session.ValidateSession(sbvision.SessionJWT(r.Header.Get("Session")))
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	// Only youtube is supported at this time, here is the "polymorphic dispatch"
	switch video.Type {
	case sbvision.YoutubeVideo:
		v, err := ctx.youtube.HandleDiscover(&video)
		if err != nil {
			fmt.Println("YoutubeHandler download error", err.Error())
			http.Error(w, "Could not download video", 400)
			return
		}
		data, err := json.Marshal(v)
		if err != nil {
			fmt.Println("Failed to marshal video", err)
			http.Error(w, "Could not parse video", 500)
			return
		}
		w.Write(data)

	default:
		http.Error(w, "Unknown video type", 400)
	}
}

func (ctx *serverContext) handleVideoPage(w http.ResponseWriter, r *http.Request) {
	var videos []sbvision.Video
	var total int64
	ids, err := getIDs(r, []string{"offset", "count"})
	if err != nil {
		ids, err := getIDs(r, []string{"id"})
		if err != nil {
			http.Error(w, "Please pass an id, or and offset and count", 400)
			return
		}

		videoID := ids[0]
		video, err := ctx.db.GetVideoByID(videoID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "not found", 404)
			return
		}
		videos = []sbvision.Video{*video}
		total = 1

	} else {
		offset, count := ids[0], ids[1]

		videos, err = ctx.db.GetVideos(offset, count)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error listing videos", 500)
			return
		}

		total, err = ctx.db.GetVideoCount()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error enumerating videos", 500)
			return
		}
	}

	// wrap the list in a json object
	data, err := json.Marshal(&struct {
		Videos []sbvision.Video `json:"videos"`
		Total  int64            `json:"total"`
	}{
		Videos: videos,
		Total:  total,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not get video list", 500)
		return
	}

	w.Write(data)
}
