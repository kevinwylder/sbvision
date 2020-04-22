package video

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kevinwylder/sbvision"
)

// FindVideoSource decides which source should be used for the given uri, and makes sure it's valid
func FindVideoSource(uri string) (*os.File, string, sbvision.VideoType, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, "", 0, err
	}

	var url, title string
	var video sbvision.VideoType

	switch u.Host {

	case "www.reddit.com":
		if !strings.HasPrefix(u.Path, "/r/skateboarding/comments/") {
			return nil, "", 0, fmt.Errorf("reddit post must be a /r/skateboarding comments link")
		}

		comments := "https://www.reddit.com" + strings.Join(strings.Split(u.Path, "/")[0:5], "/") + ".json"
		reddit, err := GetRedditPost(comments)
		if err != nil {
			return nil, "", 0, err
		}
		url = reddit.Media.RedditVideo.URL
		title = reddit.PostTitle
		video = sbvision.RedditVideo

	case "www.youtube.com":
		fallthrough
	case "youtu.be":
		yt, err := GetYoutubeDl(uri)
		if err != nil {
			return nil, "", 0, err
		}
		url = yt.URL()
		title = yt.PostTitle
		video = sbvision.YoutubeVideo

	default:
		return nil, "", 0, fmt.Errorf("Unknown Host: %s", u.Host)
	}

	file, err := download(url)
	if err != nil {
		return nil, "", 0, err
	}

	return file, title, video, nil
}

// downloads the file and returns the open handle
func download(url string) (*os.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}

	n, err := io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return nil, err
	}
	fmt.Println("Downloaded", n, "bytes from", url)
	file.Seek(0, 0)
	return file, nil
}
