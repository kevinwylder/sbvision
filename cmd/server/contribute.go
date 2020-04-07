package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/sbimage"
)

func (ctx *serverContext) handleAddFrame(w http.ResponseWriter, r *http.Request) {

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

	framePage, err := ctx.db.DataWhereFrame(hash, frameNum, video)
	if err != nil {
		fmt.Println(err)
	}
	var frame *sbvision.Frame

	if framePage != nil && len(framePage.Frames) > 0 {
		frame = &framePage.Frames[0]
	} else {
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

		err = ctx.assets.PutFrame(frame.ID, file)
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

func (ctx *serverContext) handleAddBounds(w http.ResponseWriter, r *http.Request) {
	frameid, err := getIDs(r, []string{"frame"})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session, err := ctx.session.ValidateSession(sbvision.SessionJWT(r.Header.Get("Session")))
	if err != nil {
		http.Error(w, "Missing session token", 401)
		return
	}

	frame := sbvision.Frame{
		ID: frameid[0],
		Bounds: []sbvision.Bound{
			sbvision.Bound{
				FrameID: frameid[0],
			},
		},
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&frame.Bounds[0])
	if err != nil {
		http.Error(w, "could not parse body", 400)
		return
	}

	err = ctx.db.AddBounds(&frame.Bounds[0], session)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not store bounds", 500)
		return
	}

	go func() {
		image, err := ctx.assets.GetFrame(frameid[0])
		if err != nil {
			return
		}

		defer image.Close()
		cropped, err := sbimage.Crop(image)
		if err != nil {
			fmt.Println(err)
			return
		}

		buffer := new(bytes.Buffer)
		err = cropped.GetCroppedPng(&frame.Bounds[0], buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = ctx.assets.PutBound(frame.Bounds[0].ID, buffer)
		if err != nil {
			fmt.Println(err)
		}
	}()

	data, err := json.Marshal(&frame.Bounds[0])
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not format bounds", 500)
		return
	}

	w.Write(data)

}

func (ctx *serverContext) handleAddRotation(w http.ResponseWriter, r *http.Request) {
	boundID, err := getIDs(r, []string{"bound"})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session, err := ctx.session.ValidateSession(sbvision.SessionJWT(r.Header.Get("Session")))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var rotation sbvision.Rotation
	err = decoder.Decode(&rotation)
	if err != nil {
		http.Error(w, "Could not parse json", 400)
		return
	}
	rotation.BoundID = boundID[0]
	err = ctx.db.AddRotation(&rotation, session)
	if err != nil {
		fmt.Println("Error storing rotation", err)
		http.Error(w, "Error storing rotation", 500)
		return
	}

	data, err := json.Marshal(rotation)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not marshal json", 500)
		return
	}

	w.Write(data)
}
