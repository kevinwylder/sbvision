package main

import (
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
