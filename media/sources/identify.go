package sources

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kevinwylder/sbvision"
)

// FindVideoSource decides which source should be used for the given uri, and makes sure it's valid
func FindVideoSource(uri string) (sbvision.VideoSource, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Host {

	case "www.reddit.com":
		if strings.HasPrefix(u.Path, "/r/skateboarding/comments/") {
			comments := "https://www.reddit.com" + strings.Join(strings.Split(u.Path, "/")[0:5], "/") + ".json"
			return GetRedditPost(comments)
		}
		return nil, fmt.Errorf("reddit post must be a /r/skateboarding comments link")

	case "www.youtube.com":
		fallthrough
	case "youtu.be":
		return GetYoutubeDl(uri)

	}

	return nil, fmt.Errorf("Unknown Host: %s", u.Host)
}
