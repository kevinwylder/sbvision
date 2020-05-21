package main

import (
	"log"
	"net/http"
	"strings"
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

	case "/video/info":
		ctx.handleGetVideo(w, r)

	case "/video/list":
		ctx.handleVideoPage(w, r)

	case "/video/upload":
		ctx.handleVideoDiscovery(w, r)

	case "/video/status":
		ctx.handleVideoStatus(w, r)

	case "/clip/list":
		ctx.handleGetClips(w, r)

	case "/clip/info":
		ctx.handleGetClip(w, r)

	case "/clip/upload":
		ctx.handleAddClip(w, r)

	case "/user":
		ctx.handleGetUserInfo(w, r)

	case "/remote/desktop":
		ctx.handleDesktopConnection(w, r)

	case "/remote/phone":
		ctx.handlePhoneConnection(w, r)

	default:
		if strings.HasPrefix(r.URL.Path, "/sns") {
			ctx.processes.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Not Found", 404)

	}

}
