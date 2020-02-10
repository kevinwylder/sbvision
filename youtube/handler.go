package youtube

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/kevinwylder/sbvision"
)

// YTDatabase is the required subset of functions used for this
type YTDatabase interface {
	sbvision.ImageTracker
	sbvision.VideoTracker
	sbvision.YoutubeVideoTracker
	sbvision.YoutubeSearch
}

type youtubeHandler struct {
	db     YTDatabase
	images sbvision.KeyValueStore
	proxy  *httputil.ReverseProxy
	// cache maps the video id to it's youtube info
	cache map[int64]*sbvision.YoutubeVideoInfo
}

// NewYoutubeHandler constructs a namespace for downloading youtube video info
func NewYoutubeHandler(storage YTDatabase, images sbvision.KeyValueStore) sbvision.VideoHandler {
	handler := &youtubeHandler{
		db:     storage,
		images: images,
		cache:  make(map[int64]*sbvision.YoutubeVideoInfo),
	}
	handler.proxy = &httputil.ReverseProxy{
		Director: handler.director,
	}
	return handler
}

func (dl *youtubeHandler) HandleDiscover(req *sbvision.VideoDiscoverRequest) (*sbvision.Video, error) {
	if req.Type != 1 {
		return nil, fmt.Errorf("\n\tDownload request is not a youtube type")
	}
	if req.Session == nil {
		return nil, fmt.Errorf("\n\tMissing Session from videodownloadrequest")
	}
	yt, video, err := dl.getYoutubeVideo(req.URL)
	if err != nil {
		return nil, fmt.Errorf("\n\tCould not get video: %s", err.Error())
	}
	// add video and thumbnail image to database
	err = dl.db.AddImage(video.Thumbnail, req.Session)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot add video thumbnail: %s", err.Error())
	}
	err = dl.db.AddVideo(video)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot add video: %s", err.Error())
	}
	yt.VideoID = video.ID
	err = dl.db.AddYoutubeRecord(yt)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot add youtube video: %s", err.Error())
	}
	return video, nil
}

func (dl *youtubeHandler) HandleStream(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "video id must be an int", 400)
		return
	}
	var video *sbvision.YoutubeVideoInfo
	var exists bool
	if video, exists = dl.cache[id]; !exists {
		video, err = dl.db.GetYoutubeRecord(id)
		if err != nil {
			http.Error(w, "Could not find video", 404)
			return
		}
		dl.cache[id] = video
	}
	if time.Now().After(video.MirrorExp) {
		err = dl.updateVideoLink(video)
		if err != nil {
			http.Error(w, "Could not refresh video link", 500)
			return
		}
		err = dl.db.AddYoutubeRecord(video)
		if err != nil {
			fmt.Printf("Error updating youtube mirror: %s\n", err.Error())
		}
	}
	dl.proxy.ServeHTTP(w, r)
}

func (dl *youtubeHandler) director(r *http.Request) {
	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	if err != nil {
		fmt.Println("youtube proxy directory couldn't parse id:", err.Error())
		return
	}
	if _, exists := dl.cache[id]; !exists {
		fmt.Println("youtube proxy director could not find id in cache: ", id)
		return
	}
	url, err := url.Parse(dl.cache[id].MirrorURL)
	if err != nil {
		fmt.Println("youtube proxy director couldn't parse mirror url", url)
		return
	}
	r.URL = url
}
