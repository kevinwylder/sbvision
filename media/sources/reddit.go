package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

type redditPage struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string     `json:"kind"`
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// RedditPost are reddit comments as they come from the website
type RedditPost struct {
	url       string
	ID        string `json:"id"`
	PostTitle string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Media     struct {
		RedditVideo struct {
			Duration int64  `json:"duration"`
			URL      string `json:"fallback_url"`
			IsGIF    bool   `json:"is_gif"`
		} `json:"reddit_video"`
	} `json:"media"`
}

// GetRedditSkateboardingPosts reads the frontpage of /r/skateboarding and returns all the comments links
func GetRedditSkateboardingPosts() ([]string, error) {
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/skateboarding.json", nil)
	if err != nil {
		return nil, fmt.Errorf("\n\tError opening /r/skateboarding: %s", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x86; rv:59.0) Gecko/20100101 Firefox/59.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("\n\tError doing /r/skateboarding request: %s", err)
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var page redditPage
	err = decoder.Decode(&page)
	if err != nil {
		return nil, fmt.Errorf("\n\tError decoding body: %s", err)
	}

	var posts []string
	for i := range page.Data.Children {
		url := fmt.Sprintf("https://www.reddit.com/r/skateboarding/comments/%s.json", page.Data.Children[i].Data.ID)
		posts = append([]string{url}, posts...)
	}
	return posts, nil
}

// GetRedditPost reads the url of the reddit comments and gets the json info
func GetRedditPost(url string) (*RedditPost, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("\n\tError opening reddit web request: %s", err.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("\n\tError doing reddit web request: %s", err.Error())
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var data []redditPage
	err = decoder.Decode(&data)

	// check for formatting
	if err != nil {
		return nil, fmt.Errorf("\n\tError reading reddit comments: %s", err.Error())
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("\n\tNo data found")
	}
	if len(data[0].Data.Children) == 0 {
		return nil, fmt.Errorf("\n\tEmpty top level comment")
	}
	if data[0].Data.Children[0].Data.Media.RedditVideo.IsGIF {
		// return nil, fmt.Errorf("\n\tUnsupported type: GIF")
	}
	data[0].Data.Children[0].Data.url = url
	return &data[0].Data.Children[0].Data, nil
}

// GetVideo is a constructor for the video type
func (info *RedditPost) GetVideo() sbvision.Video {
	return sbvision.Video{
		Title:     info.PostTitle,
		Type:      sbvision.RedditVideo,
		SourceURL: info.Media.RedditVideo.URL,
		ShareURL:  info.url,
	}
}

// Cleanup does nothing
func (info *RedditPost) Cleanup() {
}

// GetThumbnail gets the thumbnail from this posts and stores it in the key value store
func (info *RedditPost) GetThumbnail() (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", info.Thumbnail, nil)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting thumbnail: %s", err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("\n\tError doing thumbnail request: %s", err.Error())
	}
	return res.Body, nil
}
