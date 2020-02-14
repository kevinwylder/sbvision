package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

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

	rotation := &sbvision.Rotation{}
	err = decoder.Decode(rotation)
	if err != nil {
		http.Error(w, "Could not parse json", 400)
		return
	}
	rotation.BoundID = boundID[0]

	err = ctx.db.AddRotation(rotation, session)
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

func (ctx *serverContext) handleGetRotationFrames(w http.ResponseWriter, r *http.Request) {
	frames, err := ctx.db.DataRotationFrames()
	if err != nil {
		fmt.Println("Error getting rotation frames:", err)
		http.Error(w, "Could not get frames", 500)
		return
	}

	data, err := json.Marshal(&struct {
		Frames []sbvision.Frame `json:"frames"`
	}{
		Frames: frames,
	})
	if err != nil {
		fmt.Println("Error formatting rotation frames:", err)
		http.Error(w, "Could not get frames", 500)
		return
	}
	w.Write(data)

}
