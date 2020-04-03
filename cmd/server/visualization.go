package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database"

	"github.com/gorilla/websocket"
)

type visualizor struct {
	stopped bool
	conn    *websocket.Conn
	assets  sbvision.KeyValueStore
	db      *database.SBDatabase
	cache   databaseCache

	rchan    chan struct{}
	rmutex   sync.Mutex
	rotation sbvision.Rotation
	fchan    chan struct{}
	fmutex   sync.Mutex
	frame    sbvision.Frame
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

		v.rmutex.Lock()
		v.rotation = localRotation
		v.rmutex.Unlock()

		if v.rchan != nil {
			t := v.rchan
			v.rchan = nil
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

		if v.rotation.R == lookup.R &&
			v.rotation.I == lookup.I &&
			v.rotation.J == lookup.J &&
			v.rotation.K == lookup.K {

			v.rchan = make(chan struct{})
			<-v.rchan
			continue
		}

		v.rmutex.Lock()
		lookup = v.rotation
		v.rmutex.Unlock()

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

		v.fmutex.Lock()
		v.frame = sbvision.Frame{
			Bounds: []sbvision.Bound{
				sbvision.Bound{
					ID: nearest.BoundID,
					Rotations: []sbvision.Rotation{
						*nearest,
					},
				},
			},
		}
		v.fmutex.Unlock()

		if v.fchan != nil {
			t := v.fchan
			v.fchan = nil
			close(t)
		}
	}
}

func (v *visualizor) write() {
	var key sbvision.Key
	var id int64
	for {
		if v.stopped {
			return
		}

		if v.frame.Bounds == nil || id == v.frame.Bounds[0].ID {
			v.fchan = make(chan struct{})
			<-v.fchan
			continue
		}

		v.fmutex.Lock()
		key = v.frame.Bounds[0].Key()
		id = v.frame.Bounds[0].ID
		v.fmutex.Unlock()

		reader, err := v.assets.GetAsset(key)
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

		r := v.frame.Bounds[0].Rotations[0]
		_, err = writer.Write([]byte(fmt.Sprintf(`{"r":[%f,%f,%f,%f],"s":"data:image/png;base64,`, r.R, r.I, r.J, r.K)))
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
