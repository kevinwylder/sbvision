package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleGetFrames(w http.ResponseWriter, r *http.Request) {
	var offset int64
	if ids, err := getIDs(r, []string{"offset"}); err == nil {
		offset = ids[0]
	}

	ids, err := getIDs(r, []string{"video"})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	page, err := ctx.db.DataWhereVideo(ids[0], offset)

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

func (ctx *serverContext) handleAPIImage(w http.ResponseWriter, r *http.Request) {
	var image io.ReadCloser
	var err error
	dispatchErr := urlParamDispatch(r.Form, []idDispatch{
		idDispatch{
			keys:        []string{"bound"},
			description: "The cropped image for the given bounds id",
			handler: func(ids []int64) {
				bound := &sbvision.Bound{
					ID: ids[0],
				}
				image, err = ctx.assets.GetAsset(bound.Key())
			},
		},
		idDispatch{
			keys:        []string{"frame"},
			description: "The whole image for the given frame ID",
			handler: func(ids []int64) {
				frame := &sbvision.Frame{
					ID: ids[0],
				}
				image, err = ctx.assets.GetAsset(frame.Key())
			},
		},
	})

	if dispatchErr != nil {
		http.Error(w, dispatchErr.Error(), 400)
		return
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "image not found", 404)
		return
	}

	defer image.Close()

	data, err := ioutil.ReadAll(image)
	if err != nil {
		fmt.Println("API Get image error: ", err)
		http.Error(w, "could not read image", 500)
		return
	}

	w.Write(data)
}
