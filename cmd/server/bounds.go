package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleBoundsUpload(w http.ResponseWriter, r *http.Request) {
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

	bounds := sbvision.Bound{
		FrameID: frameid[0],
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&bounds)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not parse body", 400)
		return
	}

	err = ctx.db.AddBounds(&bounds, session)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not store bounds", 500)
		return
	}

	data, err := json.Marshal(&bounds)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not format bounds", 500)
		return
	}

	w.Write(data)

}

func (ctx *serverContext) handleBoundsImage(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, []string{"id"})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	frame, err := ctx.db.DataByBoundID(ids[0])
	if err != nil {
		http.Error(w, "Could not find bounds", 404)
		return
	}

	image, err := ctx.cropper.GetFrame(frame)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not decode image", 500)
		return
	}

	buffer := new(bytes.Buffer)
	err = image.GetCroppedPng(&frame.Bounds[0], buffer)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not crop png", 500)
		return
	}

	w.Write(buffer.Bytes())
}
