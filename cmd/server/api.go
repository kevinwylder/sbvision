package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleGetFrames(w http.ResponseWriter, r *http.Request) {
	var offset int64
	if ids, err := getIDs(r, []string{"offset"}); err == nil {
		offset = ids[0]
	}

	var err error
	var page *sbvision.FramePage

	dispatchErr := urlParamDispatch(r.Form, []idDispatch{
		idDispatch{
			description: "all frames for the given video",
			keys:        []string{"video"},
			handler: func(ids []int64) {
				page, err = ctx.db.DataWhereVideo(ids[0], offset)
			},
		},
		idDispatch{
			description: "all frames that have no rotation",
			keys:        []string{"rotationless"},
			handler: func(ids []int64) {
				page, err = ctx.db.DataWhereNoRotation(offset)
			},
		},
	})

	if dispatchErr != nil {
		http.Error(w, dispatchErr.Error(), 400)
		return
	}

	if err != nil {
		fmt.Println("Error in API endpoint", err)
		http.Error(w, "An error occurred", 500)
		return
	}

	data, err := json.Marshal(page)
	if err != nil {
		fmt.Println("Error formatting video frame data", err)
		http.Error(w, "Could not format response", 500)
		return
	}

	w.Write(data)
}
