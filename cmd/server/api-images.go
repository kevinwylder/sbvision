package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

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
		http.Error(w, err.Error(), 400)
		return
	}

	if err != nil {
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
