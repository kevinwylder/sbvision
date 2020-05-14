package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/video"
)

func (ctx *serverContext) handleGetClip(w http.ResponseWriter, r *http.Request) {
	if r.Form.Get("id") == "" {
		http.Error(w, "Missing ?id= query param", 400)
		return
	}

	clip, err := ctx.ddb.GetClipByID(r.Form.Get("id"))
	if err != nil {
		http.Error(w, "Could not get that id", 404)
		return
	}

	json.NewEncoder(w).Encode(clip)
}

func (ctx *serverContext) handleGetClips(w http.ResponseWriter, r *http.Request) {

	clips, err := ctx.ddb.GetClips(r.Form.Get("trick"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not get clips", 500)
		return
	}

	json.NewEncoder(w).Encode(clips)

}

func (ctx *serverContext) handleAddClip(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.auth.User(r.Header.Get("Identity"))
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	var clip sbvision.Clip
	err = json.NewDecoder(r.Body).Decode(&clip)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse request body", 400)
		return
	}

	go ctx.startBatchProcess(&clip)

	err = ctx.ddb.AddClip(&clip, user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add clip to database", 500)
		return
	}

	json.NewEncoder(w).Encode(&clip)
}

func (ctx *serverContext) startBatchProcess(clip *sbvision.Clip) {
	data, _ := json.Marshal(clip)
	ctx.batch.SubmitJob(&batch.SubmitJobInput{
		JobDefinition: aws.String("sbclipvid"),
		JobQueue:      aws.String(video.BatchQueueName),
		JobName:       aws.String(clip.ID),
		ContainerOverrides: &batch.ContainerOverrides{
			Command: []*string{
				aws.String("sbclipvid"),
				aws.String(base64.URLEncoding.EncodeToString(data)),
			},
		},
	})
}
