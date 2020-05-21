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
	phone   *websocket.Conn
	desktop *websocket.Conn

	Active     bool                `json:"active"`
	Locked     bool                `json:"locked"`
	LockedOn   sbvision.Quaternion `json:"lockedOn"`
	Correction sbvision.Quaternion `json:"correction"`
}

// route to forward readings from phone to desktop
func (ctx *serverContext) handlePhoneConnection(w http.ResponseWriter, r *http.Request) {
	// login
	user, err := ctx.auth.User(r.Form.Get("identity"))
	if err != nil {
		http.Error(w, "Missing identity token", 401)
		return
	}

	// connect
	socket, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("phone connect", err)
		return
	}

	// construct
	session := ctx.remotes[user.Email]
	if session == nil {
		session = &remoteSession{}
		ctx.remotes[user.Email] = session
	}
	session.connectPhone(socket)

	// forwarding
	errCounter := 0
	for {
		// read from phone
		var reading message
		err := socket.ReadJSON(&reading)
		if err != nil {
			fmt.Println("phone read", err)
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
				fmt.Println("phone write", err)
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
	session.disconnectPhone(socket)
}

func (s *remoteSession) connectPhone(socket *websocket.Conn) {
	s.mutex.Lock()
	if s.phone != socket {
		s.phone.Close()
	}
	s.phone = socket
	s.sync()
	s.mutex.Unlock()
}

func (s *remoteSession) disconnectPhone(socket *websocket.Conn) {
	if s.phone != socket {
		return
	}
	s.mutex.Lock()
	s.phone.Close()
	s.phone = nil
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
		session = &remoteSession{}
		ctx.remotes[user.Email] = session
	}
	session.connectDesktop(socket)

	// forwarding
	for {
		// read from desktop
		var read message
		err := socket.ReadJSON(&read)
		if err != nil {
			fmt.Println("desktop read", err)
			break
		}

		// write to phone
		session.mutex.Lock()
		switch read.Command {
		case "lock":
			session.Locked = true
			session.LockedOn = read.Quaternion
		case "correct":
			session.Correction = read.Quaternion
		}
		session.sync()
		session.mutex.Unlock()

	}
	session.disconnectDesktop(socket)
}

func (s *remoteSession) connectDesktop(socket *websocket.Conn) {
	s.mutex.Lock()
	if s.desktop != socket {
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
	s.Active = (s.phone != nil) && (s.desktop != nil)
	fmt.Println("sync")
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(s)
	if s.phone != nil {
		err := s.phone.WriteJSON(s)
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
