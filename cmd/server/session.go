package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) getSession(w http.ResponseWriter, r *http.Request) {
	session := &sbvision.Session{
		Time: time.Now().Unix(),
	}
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		session.IP = forwarded
	} else {
		session.IP = r.RemoteAddr
	}
	err := ctx.db.AddSession(session)
	if err != nil {
		fmt.Println("Create session error: ", err.Error())
		http.Error(w, "Could not create session", 500)
		return
	}

	jwt, err := ctx.session.SignSession(session)
	if err != nil {
		fmt.Println("Sign session error: ", err)
		http.Error(w, "Could not create session", 500)
		return
	}

	w.Write([]byte(jwt))
}
