package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database"
	"github.com/kevinwylder/sbvision/media"

	"github.com/gorilla/websocket"
)

type visualizor struct {
	stopped bool
	conn    *websocket.Conn
	assets  *media.AssetDirectory
	db      *database.SBDatabase
	cache   databaseCache

	ingoing sbvision.Rotation
	inchan  chan struct{}
	inmutex sync.Mutex

	outgoing sbvision.Rotation
	outchan  chan struct{}
	outmutex sync.Mutex
}

func (ctx *serverContext) handleVisualizationSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ctx.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
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
	var localRotation sbvision.Rotation
	for {
		if v.stopped {
			return
		}

		err := v.conn.ReadJSON(&localRotation)

		if err != nil {
			v.stopped = true
			v.conn.Close()
			return
		}

		v.inmutex.Lock()
		v.ingoing = localRotation
		v.inmutex.Unlock()

		if v.inchan != nil {
			t := v.inchan
			v.inchan = nil
			close(t)
		}
	}
}

func (v *visualizor) lookup() {
	var lookup sbvision.Rotation
	v.cache.minProduct = 2
	for {
		if v.stopped {
			return
		}

		if v.ingoing.R == lookup.R &&
			v.ingoing.I == lookup.I &&
			v.ingoing.J == lookup.J &&
			v.ingoing.K == lookup.K {

			v.inchan = make(chan struct{})
			<-v.inchan
			continue
		}

		v.inmutex.Lock()
		lookup = v.ingoing
		v.inmutex.Unlock()

		nearest := v.cache.check(&lookup)
		if nearest == nil {
			err := v.cache.load(v.db, &lookup)
			if err != nil {
				fmt.Println(err)
				v.stopped = true
				v.conn.Close()
				return
			}
			nearest = v.cache.check(&lookup)
			if nearest == nil {
				fmt.Println("Could not find rotation in cache after refresh!")
				continue
			}
		}

		v.outmutex.Lock()
		v.outgoing = *nearest
		v.outmutex.Unlock()

		if v.outchan != nil {
			t := v.outchan
			v.outchan = nil
			close(t)
		}
	}
}

func (v *visualizor) write() {
	var rotation sbvision.Rotation
	for {
		if v.stopped {
			return
		}

		if rotation == v.outgoing {
			v.outchan = make(chan struct{})
			<-v.outchan
			continue
		}

		v.outmutex.Lock()
		rotation = v.outgoing
		v.outmutex.Unlock()

		reader, err := v.assets.GetBound(rotation.BoundID)
		if err != nil {
			fmt.Println("Asset error", err)
			continue
		}

		if v.stopped {
			reader.Close()
			return
		}
		writer, err := v.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			v.stopped = true
			reader.Close()
			v.conn.Close()
			return
		}

		_, err = writer.Write([]byte(fmt.Sprintf(`{"r":[%f,%f,%f,%f],"s":"data:image/png;base64,`, rotation.R, rotation.I, rotation.J, rotation.K)))
		if err != nil {
			fmt.Println("Write err", err)
			reader.Close()
			writer.Close()
			continue
		}

		b64writer := base64.NewEncoder(base64.StdEncoding, writer)
		_, err = io.Copy(b64writer, reader)
		b64writer.Close()
		reader.Close()
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

type databaseCache struct {
	rotations  [100]sbvision.Rotation
	search     sbvision.Rotation
	minProduct float64
}

func dot(a, b *sbvision.Rotation) float64 {
	return a.R*b.R + a.I*b.I + a.J*b.J + a.K*b.K
}

func (cache *databaseCache) load(db *database.SBDatabase, rotation *sbvision.Rotation) error {
	err := db.DataNearestRotation(rotation, cache.rotations[:])
	if err != nil {
		return err
	}
	cache.search = *rotation
	cache.minProduct = dot(&cache.rotations[99], rotation)
	return nil
}

func (cache *databaseCache) check(search *sbvision.Rotation) *sbvision.Rotation {
	if dot(search, &cache.search) < cache.minProduct {
		return nil
	}
	var closest *sbvision.Rotation
	var best float64 = -1
	for i := range cache.rotations {
		product := dot(search, &cache.rotations[i])
		if product > best {
			best = product
			closest = &cache.rotations[i]
		}
	}
	return closest
}
