package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

func (ctx *serverContext) handleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	email, err := ctx.auth.Verify(r.Header.Get("Identity"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", 401)
		return
	}

	var user *sbvision.User

	switch r.Method {
	case http.MethodPost:

		user = &sbvision.User{
			Email: email,
		}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			http.Error(w, "Missing email object in body", 400)
			return
		}
		err = ctx.db.AddUser(user)
		if err != nil {
			http.Error(w, "Could not add user to database", 500)
			return
		}

	case http.MethodGet:
		user, err = ctx.db.GetUser(email)

	}

	json.NewEncoder(w).Encode(user)
}
