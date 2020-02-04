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
		// Route to download video
		// Requires session
		err := ctx.sessionManager.ValidateSession(r.Header.Get("Session"))
		if err != nil {
			http.Error(w, "unauthorized", 401)
			return
		}

		var video sbvision.VideoDownloadRequest
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&video)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "invalid video format", 400)
			return
		}

		// Only youtube is supported at this time, here is the "polymorphic dispatch"
		switch video.Type {
		case sbvision.YoutubeVideo:
			err = ctx.youtubeDownloader.Handle(&video)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Could not download video", 400)
				return
			}

		default:
			http.Error(w, "Unknown video type", 400)
		}

	case http.MethodGet:
		// Route to get page of videos
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "invalid form", 400)
			return
		}
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

		videos, err := ctx.videoLister.GetVideos(offset, count)
		if err != nil {
			http.Error(w, "Error listing videos", 500)
			return
		}

		// wrap the list in a json object
		data, err := json.Marshal(&struct {
			Videos []sbvision.Video `json:"videos"`
		}{
			Videos: videos,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not marshal videos", 500)
			return
		}

		w.Write(data)
	}
}
