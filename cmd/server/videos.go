package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kevinwylder/sbvision/media/sources"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleVideoPage(w http.ResponseWriter, r *http.Request) {
	var videos []sbvision.Video
	var total int64
	var err error

	dispatchErr := urlParamDispatch(r.Form, []idDispatch{
		idDispatch{
			description: "a page of video results",
			keys:        []string{"offset", "count"},
			handler: func(ids []int64) {
				offset, count := ids[0], ids[1]

				videos, err = ctx.db.GetVideos(offset, count)
				if err != nil {
					http.Error(w, "Error listing videos", 500)
					return
				}

				total, err = ctx.db.GetVideoCount()
				if err != nil {
					http.Error(w, "Error enumerating videos", 500)
					return
				}
			},
		},
		idDispatch{
			description: "A single video",
			keys:        []string{"id"},
			handler: func(ids []int64) {
				videoID := ids[0]
				video, err := ctx.db.GetVideoByID(videoID)
				if err != nil {
					http.Error(w, "not found", 404)
					return
				}
				videos = append(videos, *video)
				total = 1
			},
		},
	})

	if dispatchErr != nil {
		http.Error(w, dispatchErr.Error(), 400)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	// wrap the list in a json object
	json.NewEncoder(w).Encode(&struct {
		Videos []sbvision.Video `json:"videos"`
		Total  int64            `json:"total"`
	}{
		Videos: videos,
		Total:  total,
	})
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

	source, err := sources.FindVideoSource(request.URL)
	if err != nil {
		http.Error(w, "Error getting video: "+err.Error(), 400)
		return
	}

	ticket, err := ctx.discoveryQueue.Enqueue(user, source)
	if err != nil {
		http.Error(w, "Queue is full, come back later", 503)
		return
	}
	json.NewEncoder(w).Encode(ticket)
}

func (ctx *serverContext) handleVideoStatus(w http.ResponseWriter, r *http.Request) {

}
