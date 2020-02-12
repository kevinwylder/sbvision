package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleFrameUpload(w http.ResponseWriter, r *http.Request) {

	session, err := ctx.session.ValidateSession(sbvision.SessionJWT(r.Header.Get("Session")))
	if err != nil {
		http.Error(w, "Missing session token", 401)
		return
	}

	ids, err := getIDs(r, []string{"video", "frame"})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	video, frameNum := ids[0], ids[1]

	frame, err := ctx.db.GetFrame(video, frameNum)
	if frame == nil {

		err := r.ParseMultipartForm(10 * 1024 * 1024)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error parsing multipart form", 400)
			return
		}

		image := sbvision.Image(fmt.Sprintf("frame/%d-%d.png", video, frameNum))
		file, _, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Could not get image file", 400)
			return
		}

		err = ctx.assets.PutAsset(string(image), file)
		if err != nil {
			fmt.Println("Error putting asset", err)
			http.Error(w, "Error storing asset", 500)
			return
		}

		err = ctx.db.AddImage(image, session)
		if err != nil {
			http.Error(w, "Error adding image to DB", 500)
			return
		}

		frame = &sbvision.Frame{
			VideoID: video,
			Time:    frameNum,
			Image:   image,
		}

		err = ctx.db.AddFrame(frame)
		if err != nil {
			fmt.Println("Error adding frame", err)
			http.Error(w, "Could not add frame", 500)
			return
		}

	}

	data, err := json.Marshal(frame)
	if err != nil {
		http.Error(w, "Error representing frame", 500)
		return
	}

	w.Write(data)
}

func (ctx *serverContext) handleGetFrames(w http.ResponseWriter, r *http.Request) {
	video, err := getIDs(r, []string{"video"})
	if err != nil {
		frames, bounds, rotations, err := ctx.db.DataCounts()
		if err != nil {
			http.Error(w, "Error counting data", 500)
			return
		}
		data, err := json.Marshal(&struct {
			Frames    int64 `json:"frames"`
			Bounds    int64 `json:"bounds"`
			Rotations int64 `json:"rotations"`
		}{
			Frames:    frames,
			Bounds:    bounds,
			Rotations: rotations,
		})
		if err != nil {
			http.Error(w, "Error formating counted data", 500)
			return
		}
		w.Write(data)
		return
	}
	videoID := video[0]
	frames, err := ctx.db.DataVideoFrames(videoID)
	if err != nil {
		fmt.Println("error loading video frames from db", err)
		http.Error(w, "Error loading video frames", 500)
		return
	}

	data, err := json.Marshal(&struct {
		Frames []sbvision.Frame `json:"frames"`
	}{
		Frames: frames,
	})
	if err != nil {
		fmt.Println("Error formatting video frame data", err)
		http.Error(w, "Could not format response", 500)
		return
	}

	w.Write(data)
}
