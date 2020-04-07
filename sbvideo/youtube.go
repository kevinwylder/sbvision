package sbvideo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kevinwylder/sbvision"
)

// YoutubeDl is the information youtube-dl 's .info.json file structure
type YoutubeDl struct {
	Thumbnail string `json:"thumbnail"`
	Title     string `json:"title"`
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
	log.Println("Downloading video info", url, "to", directory)

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
	return &info, nil
}

// Update puts this info into a video struct
func (info *YoutubeDl) Update(video *sbvision.Video) {
	video.Title = info.Title
	video.Duration = info.Duration
	video.Type = sbvision.YoutubeVideo
	video.Format = "video/mp4"

	var largestFormat int64
	for _, format := range info.Formats {
		if format.Filesize < largestFormat {
			continue
		}
		largestFormat = format.Filesize
		video.URL = format.URL
	}

	expireMatcher := regexp.MustCompile(`expire=(\d+)`)
	expires := expireMatcher.FindStringSubmatch(video.URL)
	if len(expires) == 0 {
		video.LinkExp = time.Now()
	} else {
		unix, _ := strconv.ParseInt(expires[1], 10, 64)
		video.LinkExp = time.Unix(unix, 0)
	}

}

// GetThumbnail downloads the thumbnail for the video
func (info *YoutubeDl) GetThumbnail() (io.ReadCloser, error) {
	res, err := http.Get(info.Thumbnail)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting thumbnail: %s", err.Error())
	}
	return res.Body, nil
}
