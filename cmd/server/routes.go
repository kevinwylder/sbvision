package main

import (
	"log"
	"net/http"
)

func (ctx *serverContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	switch r.URL.Path {

	case "/app/session":
		ctx.handleNewSession(w, r)

	case "/app/video/list":
		ctx.handleVideoPage(w, r)

	case "/app/video/thumbnail":
		ctx.handleVideoThumbnail(w, r)

	case "/app/video/stream":
		ctx.proxy.ServeHTTP(w, r)

	case "/app/contribute/frame":
		ctx.handleAddFrame(w, r)

	case "/app/contribute/bounds":
		ctx.handleAddBounds(w, r)

	case "/app/contribute/rotation":
		ctx.handleAddRotation(w, r)

	case "/app/visualization":
		ctx.handleVisualizationSocket(w, r)

	case "/api/frames":
		ctx.handleGetFrames(w, r)

	case "/api/image":
		ctx.handleAPIImage(w, r)

	default:
		http.Error(w, "Not Found", 404)

	}

}
