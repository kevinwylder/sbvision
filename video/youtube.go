package video

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

// YoutubeDl is the information youtube-dl 's .info.json file structure
type YoutubeDl struct {
	url       string
	Thumbnail string `json:"thumbnail"`
	PostTitle string `json:"title"`
	DisplayID string `json:"display_id"`
	Duration  int64  `json:"duration"`
	Formats   []struct {
		Filesize int64   `json:"filesize"`
		URL      string  `json:"url"`
		FPS      float64 `json:"fps"`
	} `json:"formats"`
}

// GetYoutubeDl reads the video.OriginURL and gets the information from youtube
func GetYoutubeDl(url string) (*YoutubeDl, error) {
	// create a temp dir to download the json and thumbnail
	directory, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create tmp dir: %s", err.Error())
	}
	defer os.RemoveAll(directory)

	// run the youtube-dl command to get the info
	cmd := exec.Command("youtube-dl", url, "--write-info-json", "--skip-download")
	cmd.Dir = directory
	err = cmd.Run()
	if err != nil {
		data, _ := cmd.CombinedOutput()
		return nil, fmt.Errorf("\n\tError running youtube-dl for %s %s.\nyoutube-dl Output: %s", url, err.Error(), string(data))
	}

	// look for the json file in the tmp directory
	files, err := ioutil.ReadDir(directory)
	var infoJSON os.FileInfo
	for i, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			infoJSON = files[i]
		}
	}
	if infoJSON == nil {
		return nil, fmt.Errorf("\n\tCould not find info.json")
	}

	// open and read the file
	infoFile, err := os.Open(path.Join(directory, infoJSON.Name()))
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot open .info.json file: %s", err.Error())
	}
	defer infoFile.Close()
	var info YoutubeDl
	decoder := json.NewDecoder(infoFile)
	err = decoder.Decode(&info)
	if err != nil {
		return nil, err
	}
	info.url = url
	return &info, nil
}

// URL returns the largest format video url
func (info *YoutubeDl) URL() string {
	var largestFormat int64
	var url string
	for _, format := range info.Formats {
		if format.Filesize < largestFormat {
			continue
		}
		largestFormat = format.Filesize
		url = format.URL
	}
	return url
}

// Cleanup does nothing
func (info *YoutubeDl) Cleanup() {
}

// GetThumbnail downloads the thumbnail for the video
func (info *YoutubeDl) GetThumbnail() (io.ReadCloser, error) {
	res, err := http.Get(info.Thumbnail)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting thumbnail: %s", err.Error())
	}
	return res.Body, nil
}
