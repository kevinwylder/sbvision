package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (ctx *serverContext) handleAPIImage(w http.ResponseWriter, r *http.Request) {
	var image io.ReadCloser
	var err error
	dispatchErr := urlParamDispatch(r.Form, []idDispatch{
		idDispatch{
			keys:        []string{"bound"},
			description: "The cropped image for the given bounds id",
			handler: func(ids []int64) {
				image, err = ctx.assets.GetBound(ids[0])
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
