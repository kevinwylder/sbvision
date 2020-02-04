package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) session(w http.ResponseWriter, r *http.Request) {
	session := &sbvision.Session{
		Time: time.Now().Unix(),
	}
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		session.IP = forwarded
	} else {
		session.IP = r.RemoteAddr
	}
	err := ctx.sessionStorage.TrackSession(session)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not create session", 500)
		return
	}

	jwt, err := ctx.sessionManager.SignSession(session)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not create session", 500)
		return
	}

	w.Write([]byte(jwt))
}
