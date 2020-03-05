package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database"

	"github.com/gorilla/websocket"
)

type visualizor struct {
	stopped bool
	conn    *websocket.Conn
	assets  sbvision.KeyValueStore
	db      *database.SBDatabase

	rmutex   sync.Mutex
	rotation sbvision.Rotation
	fmutex   sync.Mutex
	frame    sbvision.Frame
}

func (ctx *serverContext) handleVisualizationSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error starting websocket", 400)
		return
	}

	v := &visualizor{
		conn:   conn,
		db:     ctx.db,
		assets: ctx.assets,
	}

	go v.read()
	go v.lookup()
	go v.write()
}

func (v *visualizor) read() {
	for {
		if v.stopped {
			return
		}

		v.rmutex.Lock()
		err := v.conn.ReadJSON(&v.rotation)
		v.rmutex.Unlock()
		if err != nil {
			fmt.Println("websocket read", err)
			v.stopped = true
			v.conn.Close()
			return
		}
	}
}

func (v *visualizor) lookup() {
	var lookup sbvision.Rotation
	for {
		if v.stopped {
			return
		}

		if v.rotation.R == lookup.R &&
			v.rotation.I == lookup.I &&
			v.rotation.J == lookup.J &&
			v.rotation.K == lookup.K {

			time.Sleep(time.Millisecond * 50)
			continue
		}

		v.rmutex.Lock()
		lookup = v.rotation
		v.rmutex.Unlock()

		frame, err := v.db.DataNearestRotation(&lookup, 1)
		if err != nil {
			fmt.Println("websocket lookup", err)
			v.stopped = true
			v.conn.Close()
			return
		}

		v.fmutex.Lock()
		v.frame = *frame
		v.fmutex.Unlock()
	}
}

func (v *visualizor) write() {
	var key sbvision.Key
	var id int64
	for {
		if v.stopped {
			return
		}

		if id == v.frame.ID {
			time.Sleep(time.Millisecond * 50)
			continue
		}

		v.fmutex.Lock()
		key = v.frame.Bounds[0].Key()
		id = v.frame.ID
		v.fmutex.Unlock()

		reader, err := v.assets.GetAsset(key)
		if err != nil {
			fmt.Println("Asset error", err)
			continue
		}
		defer reader.Close()

		if v.stopped {
			return
		}
		writer, err := v.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			fmt.Println("websocket write", err)
			v.stopped = true
			v.conn.Close()
			return
		}

		r := v.frame.Bounds[0].Rotations[0]
		_, err = writer.Write([]byte(fmt.Sprintf(`{"r":[%f,%f,%f,%f],"s":"data:image/png;base64,`, r.R, r.I, r.J, r.K)))
		if err != nil {
			fmt.Println("Write err", err)
			writer.Close()
			continue
		}

		b64writer := base64.NewEncoder(base64.StdEncoding, writer)
		_, err = io.Copy(b64writer, reader)
		b64writer.Close()
		if err != nil {
			fmt.Println("Write error", err)
			writer.Close()
			continue
		}

		_, err = writer.Write([]byte(`"}`))
		writer.Close()
		if err != nil {
			fmt.Println("Write err", err)
		}
	}
}
