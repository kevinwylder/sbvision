package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision/sbimage"

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

	err = r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing multipart form", 400)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not get image file", 400)
		return
	}
	defer file.Close()

	hash, err := sbimage.HashImage(file)
	if err != nil {
		http.Error(w, "Could not compute image hash", 500)
		return
	}

	frame, err := ctx.db.GetFrame(video, frameNum, hash)
	if frame == nil {

		frame = &sbvision.Frame{
			VideoID: video,
			Time:    frameNum,
		}

		err = ctx.db.AddFrame(frame, session, hash)
		if err != nil {
			fmt.Println("Error adding frame", err)
			http.Error(w, "Could not add frame", 500)
			return
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			fmt.Println("Error seeking to beginning of uploaded file: ", err)
			http.Error(w, "Could not save frame", 500)
			return
		}

		err = ctx.assets.PutAsset(fmt.Sprintf("frame/%d.png", frame.ID), file)
		if err != nil {
			fmt.Println("Error putting asset", err)
			http.Error(w, "Error storing asset", 500)
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

func (ctx *serverContext) handleGetDataCounts(w http.ResponseWriter, r *http.Request) {
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
}

func (ctx *serverContext) handleGetFrames(w http.ResponseWriter, r *http.Request) {
	if r.Form.Get("rotationless") != "" {
		ctx.handleGetRotationFrames(w, r)
		return
	}

	video, err := getIDs(r, []string{"video"})
	if err != nil {
		ctx.handleGetDataCounts(w, r)
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
