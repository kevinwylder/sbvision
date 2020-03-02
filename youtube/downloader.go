package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kevinwylder/sbvision"
)

type dlInfo struct {
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

// calls youtube-dl to collect a thumbnail and .info.json file. The thumbnail is uploaded, and the .info.json is used to create
// a new YoutubeVideoInfo object to be tracked in the database
func (dl *youtubeHandler) getYoutubeVideo(url string) (*sbvision.Video, error) {
	// create a temp dir to download the json and thumbnail
	directory, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create tmp dir: %s", err.Error())
	}
	defer os.RemoveAll(directory)

	// run the youtube-dl command to get the info
	cmd := exec.Command("youtube-dl", url, "--write-info-json", "--write-thumbnail", "--skip-download")
	cmd.Dir = directory
	err = cmd.Run()
	if err != nil {
		data, _ := cmd.CombinedOutput()
		return nil, fmt.Errorf("\n\tError running youtube-dl for %s %s.\nyoutube-dl Output: %s", url, err.Error(), string(data))
	}

	// look for the files in the tmp directory
	files, err := ioutil.ReadDir(directory)
	var infoJSON os.FileInfo
	var thumbnail os.FileInfo
	for i, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			infoJSON = files[i]
		}
		if strings.HasSuffix(f.Name(), ".jpg") {
			thumbnail = files[i]
		}
	}
	if infoJSON == nil || thumbnail == nil {
		return nil, fmt.Errorf("\n\tCould not find thumbnail (%v) or info.json (%v)", thumbnail, infoJSON)
	}

	// parse video info json file
	video := &sbvision.Video{}
	yt := &sbvision.YoutubeVideoInfo{}
	infoFile, err := os.Open(path.Join(directory, infoJSON.Name()))
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot open .info.json file: %s", err.Error())
	}
	err = parseInfo(infoFile, yt, video)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot parse file: %s", err.Error())
	}

	// Add video to database
	err = dl.db.AddVideo(video)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot add video: %s", err.Error())
	}

	// upload the thumbnail
	thumbnailFile, err := os.Open(path.Join(directory, thumbnail.Name()))
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot open image file: %s", err.Error())
	}
	defer thumbnailFile.Close()
	err = dl.images.PutAsset(video.Thumbnail(), thumbnailFile)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot upload image: %s", err.Error())
	}

	// add youtube video to the database
	yt.VideoID = video.ID
	err = dl.db.AddYoutubeRecord(yt)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot add youtube video: %s", err.Error())
	}

	return video, nil
}

// updateVideoLink uses youtube-dl to acquire a new .info.json struct for the purposes of refreshing the
// video stream link.
func (dl *youtubeHandler) updateVideoLink(info *sbvision.YoutubeVideoInfo) error {
	// create a temp dir to download the json and thumbnail
	directory, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("\n\tCannot create tmp dir: %s", err.Error())
	}
	defer os.RemoveAll(directory)

	// run the youtube-dl command to get the info
	cmd := exec.Command("youtube-dl", "https://www.youtube.com/watch?v="+info.YoutubeID, "--write-info-json", "--skip-download")
	cmd.Dir = directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("\n\tError running youtube-dl %s", err.Error())
	}

	// look for the file in the tmp directory
	files, err := ioutil.ReadDir(directory)
	var infoJSON os.FileInfo
	for i, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			infoJSON = files[i]
		}
	}
	if infoJSON == nil {
		return fmt.Errorf("\n\tCould not find info.json (%v)", infoJSON)
	}

	infoFile, err := os.Open(path.Join(directory, infoJSON.Name()))
	if err != nil {
		return fmt.Errorf("\n\tCannot open .info.json file: %s", err.Error())
	}
	err = parseInfo(infoFile, info, nil)
	if err != nil {
		return fmt.Errorf("\n\tCannot parse .info.json file: %s", err.Error())
	}
	return nil
}

// parseInfo reads the info file from youtube-dl and extracts the info for the video/youtubevideo
func parseInfo(file *os.File, yt *sbvision.YoutubeVideoInfo, video *sbvision.Video) error {
	defer file.Close()

	var info dlInfo
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&info)
	if err != nil {
		return err
	}
	yt.YoutubeID = info.DisplayID
	if video != nil {
		video.Title = info.Title
		video.Duration = info.Duration
		video.Type = sbvision.YoutubeVideo
		video.Format = "video/mp4"
	}
	var largestFormat int64
	for _, format := range info.Formats {
		if format.Filesize < largestFormat {
			continue
		}
		largestFormat = format.Filesize
		yt.MirrorURL = format.URL
	}

	expireMatcher := regexp.MustCompile(`expire=(\d+)`)
	expires := expireMatcher.FindStringSubmatch(yt.MirrorURL)
	if len(expires) < 2 {
		return fmt.Errorf("\n\tCould not find expiration in url (%s)", yt.MirrorURL)
	}
	unix, err := strconv.ParseInt(expires[1], 10, 64)
	if err != nil {
		return fmt.Errorf("\n\tCannot parse expiration timestamp (%v)", expires)
	}
	yt.MirrorExp = time.Unix(unix, 0)
	return nil
}
