package sbvideo

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/kevinwylder/sbvision"
)

// Database is the required function to lookup and update videos
type Database interface {
	GetVideoByID(id int64) (*sbvision.Video, error)
	UpdateVideo(*sbvision.Video) error
}

// Proxy serves partial video segments for streaming
type Proxy struct {
	reverse  *httputil.ReverseProxy
	cache    map[int64]*sbvision.Video
	database Database
}

// NewVideoProxy creates a proxy handler for video streaming
func NewVideoProxy(database Database) *Proxy {
	proxy := &Proxy{
		cache:    make(map[int64]*sbvision.Video),
		database: database,
	}
	proxy.reverse = &httputil.ReverseProxy{
		Director: proxy.director,
	}
	return proxy
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "missing video ?id=", 400)
		return
	}

	if _, exists := proxy.cache[id]; !exists {
		proxy.cache[id], err = proxy.database.GetVideoByID(id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error getting video", 500)
			return
		}
	}

	if video := proxy.cache[id]; video.LinkExp.Before(time.Now()) {
		var info sbvision.VideoSource
		switch video.Type {
		case sbvision.YoutubeVideo:
			info, err = GetYoutubeDl(video.OriginURL)
		case sbvision.RedditVideo:
			info, err = GetRedditPost(video.OriginURL)
		}

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error refreshing video source", 500)
			return
		}
		info.Update(video)
		go proxy.database.UpdateVideo(video)
	}
	proxy.reverse.ServeHTTP(w, r)
}

func (proxy *Proxy) director(r *http.Request) {
	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	if err != nil {
		fmt.Println("proxy directory couldn't parse id:", err.Error())
		return
	}

	if _, exists := proxy.cache[id]; !exists {
		fmt.Println("proxy director could not find id in cache: ", id)
		return
	}

	url, err := url.Parse(proxy.cache[id].URL)
	if err != nil {
		fmt.Println("proxy director couldn't parse mirror url", url)
		return
	}

	r.Header.Set("Host", r.Host)
	r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	r.Host = url.Host
	r.URL = url
}
