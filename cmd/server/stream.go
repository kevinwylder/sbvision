package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (ctx *serverContext) handleStream(w http.ResponseWriter, r *http.Request) {
	switch r.Form.Get("type") {
	case "1":
		ctx.youtube.HandleStream(w, r)
	case "3":
		// this is a test video
		rangeHeader := r.Header.Get("Range")
		if len(rangeHeader) < 7 {
			http.Error(w, "cannot handle provided Range header", 400)
			return
		}
		startEnd := strings.Split(rangeHeader[6:], "-")
		start, err := strconv.ParseUint(startEnd[0], 10, 64)
		if err != nil {
			http.Error(w, "Cannot read starting byte", 400)
			return
		}
		end, err := strconv.ParseUint(startEnd[1], 10, 64)
		if err != nil {
			end = start + 1024*100
		}
		if end > uint64(len(ctx.testVideo)) {
			end = uint64(len(ctx.testVideo))
		}
		header := fmt.Sprintf("bytes %d-%d/%d", start, end, len(ctx.testVideo))
		w.Header().Set("Content-Range", header)
		w.WriteHeader(http.StatusPartialContent)
		w.Write(ctx.testVideo[start:end])

	default:
		http.Error(w, "Unsupported video type", 400)
	}
}
