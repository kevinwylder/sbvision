package main

import (
	"log"
	"net/http"
)

func (ctx *serverContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(200)
		return
	}

	log.Println(r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	switch r.URL.Path {

	case "/video/list":
		ctx.handleVideoPage(w, r)

	case "/video/thumbnail":
		ctx.handleVideoThumbnail(w, r)

	case "/video/stream":
		ctx.proxy.ServeHTTP(w, r)

	case "/frames":
		ctx.handleGetFrames(w, r)

	case "/image":
		ctx.handleAPIImage(w, r)

	case "/user":
		ctx.handleGetUserInfo(w, r)

	case "/visualization":
		ctx.handleVisualizationSocket(w, r)

	default:
		http.Error(w, "Not Found", 404)

	}

}
