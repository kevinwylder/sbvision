package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kevinwylder/sbvision"
)

type message struct {
	Command    string              `json:"cmd"`
	Quaternion sbvision.Quaternion `json:"rot"`
}

type remoteSession struct {
	mutex   sync.Mutex
	remote  *websocket.Conn
	desktop *websocket.Conn

	Active   bool                `json:"active"`
	Locked   bool                `json:"locked"`
	LockedOn sbvision.Quaternion `json:"lockedOn"`
}

// route to forward readings from remote to desktop
func (ctx *serverContext) handleRemoteConnection(w http.ResponseWriter, r *http.Request) {
	// login
	user, err := ctx.auth.User(r.Form.Get("identity"))
	if err != nil {
		http.Error(w, "Missing identity token", 401)
		return
	}

	// connect
	socket, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("remote connect", err)
		return
	}

	// construct
	session := ctx.remotes[user.Email]
	if session == nil {
		session = &remoteSession{
			remote: socket,
		}
		ctx.remotes[user.Email] = session
	} else {
		session.connectRemote(socket)
	}

	// forwarding
	errCounter := 0
	for {
		// read from remote
		var reading message
		err := socket.ReadJSON(&reading)
		if err != nil {
			fmt.Println("remote read", err)
			break
		}

		// write to desktop
		session.mutex.Lock()
		if session.Locked && reading.Quaternion.Dot(session.LockedOn) > .995 {
			session.Locked = false
		}
		if session.Active && (!session.Locked || reading.Command != "") {
			err = session.desktop.WriteJSON(&reading)
			if err != nil {
				fmt.Println("remote write", err)
				errCounter++
				if errCounter > 10 {
					fmt.Println("10th error, exiting")
					session.mutex.Unlock()
					break
				}
			} else {
				errCounter = 0
			}
		}
		session.mutex.Unlock()
	}
	session.disconnectRemote(socket)
}

func (s *remoteSession) connectRemote(socket *websocket.Conn) {
	s.mutex.Lock()
	if s.remote != nil {
		s.remote.Close()
	}
	s.remote = socket
	s.sync()
	s.mutex.Unlock()
}

func (s *remoteSession) disconnectRemote(socket *websocket.Conn) {
	if s.remote != socket {
		return
	}
	s.mutex.Lock()
	s.remote.Close()
	s.remote = nil
	s.sync()
	s.mutex.Unlock()
}

func (ctx *serverContext) handleDesktopConnection(w http.ResponseWriter, r *http.Request) {
	// login
	user, err := ctx.auth.User(r.Form.Get("identity"))
	if err != nil {
		http.Error(w, "Missing identity token", 401)
		return
	}

	// connect
	socket, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// construct
	session := ctx.remotes[user.Email]
	if session == nil {
		session = &remoteSession{
			desktop: socket,
		}
		ctx.remotes[user.Email] = session
	} else {
		session.connectDesktop(socket)
	}

	// forwarding
	for {
		// read from desktop
		var read message
		err := socket.ReadJSON(&read)
		if err != nil {
			fmt.Println("desktop read", err)
			break
		}

		// write to remote
		session.mutex.Lock()
		if read.Command == "" {
			session.Locked = true
			session.LockedOn = read.Quaternion
		}
		session.sync()
		session.mutex.Unlock()

	}
	session.disconnectDesktop(socket)
}

func (s *remoteSession) connectDesktop(socket *websocket.Conn) {
	s.mutex.Lock()
	if s.desktop != nil {
		s.desktop.Close()
	}
	s.desktop = socket
	s.sync()
	s.mutex.Unlock()
}

func (s *remoteSession) disconnectDesktop(socket *websocket.Conn) {
	if s.desktop != socket {
		return
	}
	s.mutex.Lock()
	s.desktop.Close()
	s.desktop = nil
	s.sync()
	s.mutex.Unlock()
}

func (s *remoteSession) sync() error {
	s.Active = (s.remote != nil) && (s.desktop != nil)
	fmt.Println("sync")
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(s)
	if s.remote != nil {
		err := s.remote.WriteJSON(s)
		if err != nil {
			return err
		}
	}
	if s.desktop != nil {
		err := s.desktop.WriteJSON(s)
		if err != nil {
			return err
		}
	}
	return nil
}
