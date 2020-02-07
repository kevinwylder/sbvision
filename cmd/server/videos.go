package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) videos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
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

	case http.MethodGet:
		r.ParseForm()
		// Route to get page of videos
		offset, err := strconv.Atoi(r.Form.Get("offset"))
		if err != nil {
			http.Error(w, "offset not a int", 400)
			return
		}
		count, err := strconv.Atoi(r.Form.Get("count"))
		if err != nil {
			http.Error(w, "count not a int", 400)
			return
		}

		videos, err := ctx.db.GetVideos(offset, count)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error listing videos", 500)
			return
		}

		total, err := ctx.db.GetVideoCount()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error enumerating videos", 500)
			return
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
}
