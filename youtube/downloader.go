package youtube

import (
	"encoding/json"
	"fmt"
	"io"
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
	} `json:"fommats"`
}

func getYoutubeVideo(url string, manager sbvision.ImageManager) (*sbvision.YoutubeVideoInfo, error) {
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
		return nil, fmt.Errorf("\n\tError running youtube-dl %s", err.Error())
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

	// parse video info into the
	video := sbvision.YoutubeVideoInfo{
		Video: &sbvision.Video{},
	}
	infoFile, err := os.Open(path.Join(directory, infoJSON.Name()))
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot open .info.json file: %s", err.Error())
	}
	defer infoFile.Close()
	err = parseInfo(infoFile, &video)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot parse file: %s", err.Error())
	}

	// upload the thumbnail
	thumbnailFile, err := os.Open(path.Join(directory, thumbnail.Name()))
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot open image file: %s", err.Error())
	}
	defer thumbnailFile.Close()
	video.Video.Thumbnail, err = manager.UploadImage(thumbnailFile, "thumbnail-"+video.YoutubeID+".jpg")
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot upload image: %s", err.Error())
	}

	return &video, nil
}

func parseInfo(data io.Reader, dst *sbvision.YoutubeVideoInfo) error {
	var info dlInfo
	decoder := json.NewDecoder(data)
	err := decoder.Decode(&info)
	if err != nil {
		return err
	}
	dst.Video.Title = info.Title
	dst.YoutubeID = info.DisplayID
	dst.Video.Duration = info.Duration
	dst.Video.Type = sbvision.YoutubeVideo

	var largestFormat int64
	for _, format := range info.Formats {
		if format.Filesize < largestFormat {
			continue
		}
		largestFormat = format.Filesize
		dst.MirrorURL = format.URL
		dst.Video.FPS = format.FPS
	}

	expireMatcher := regexp.MustCompile(`expire=(\d+)`)
	expires := expireMatcher.FindStringSubmatch(dst.MirrorURL)
	if len(expires) < 2 {
		return fmt.Errorf("\n\tCould not find expiration in url (%s)", dst.MirrorURL)
	}
	unix, err := strconv.ParseInt(expires[1], 10, 64)
	if err != nil {
		return fmt.Errorf("\n\tCannot parse expiration timestamp (%v)", expires)
	}
	dst.MirrorExp = time.Unix(unix, 0)

	return nil
}
