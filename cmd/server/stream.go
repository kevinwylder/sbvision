package main

import (
	"net/http"
)

func (ctx *serverContext) handleStream(w http.ResponseWriter, r *http.Request) {
	switch r.Form.Get("type") {
	case "1":
		ctx.youtube.HandleStream(w, r)

	default:
		http.Error(w, "Unsupported video type", 400)
	}
}
