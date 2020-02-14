package main

import (
	"log"
	"net/http"
	"strings"
)

func (ctx *serverContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/image") {
		ctx.handleImage(w, r)
		return
	}

	switch r.URL.Path {

	case "/videos":
		switch r.Method {
		case http.MethodPost:
			ctx.handleVideoDiscovery(w, r)
		case http.MethodGet:
			ctx.handleVideoPage(w, r)
		}

	case "/stream":
		ctx.handleStream(w, r)

	case "/session":
		ctx.handleNewSession(w, r)

	case "/frames":
		switch r.Method {
		case http.MethodPost:
			ctx.handleFrameUpload(w, r)
		case http.MethodGet:
			ctx.handleGetFrames(w, r)
		}

	case "/bounds":
		switch r.Method {
		case http.MethodPost:
			ctx.handleBoundsUpload(w, r)
		case http.MethodGet:
			ctx.handleBoundsImage(w, r)
		}

	case "/orientation":
		switch r.Method {
		case http.MethodPost:
			ctx.handleAddRotation(w, r)
		}

	default:
		// redirect /video/:id requests to index
		if strings.HasPrefix(r.URL.Path, "/video/") {
			r.URL.Path = "/"
		}
		if strings.HasPrefix(r.URL.Path, "/rotations") {
			r.URL.Path = "/"
		}
		// fallthrough to the frontend
		ctx.frontend.ServeHTTP(w, r)

	}

}
