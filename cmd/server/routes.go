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
			return
		case http.MethodGet:
			ctx.handleVideoPage(w, r)
			return
		}

	case "/stream":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		ctx.handleStream(w, r)
		return

	case "/session":
		ctx.handleNewSession(w, r)
		return

	case "/frame":
		switch r.Method {
		case http.MethodPost:
			ctx.handleFrameUpload(w, r)
			return
		}

	case "/bounds":
		switch r.Method {
		case http.MethodPost:
			ctx.handleBoundsUpload(w, r)
			return
		}

	}

	// redirect /video/:id requests to index
	if strings.HasPrefix(r.URL.Path, "/video/") {
		r.URL.Path = "/"
	}
	// fallthrough to the frontend
	ctx.frontend.ServeHTTP(w, r)
}