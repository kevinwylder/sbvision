package main

import (
	"encoding/base64"
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
	if err != nil {

		prefix := "data:image/png;base64,"
		dataURI := make([]byte, len(prefix))
		read, err := r.Body.Read(dataURI)
		if read != 22 || string(dataURI) != prefix {
			http.Error(w, "Not a valid png URI. What do you think you're tryin here...", 400)
			return
		}

		decoded := base64.NewDecoder(base64.StdEncoding, r.Body)

		image := sbvision.Image(fmt.Sprintf("frame/%d-%d.png", video, frameNum))
		err = ctx.assets.PutAsset(string(image), decoded)
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
			VideoID:  video,
			FrameNum: frameNum,
			Image:    image,
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
