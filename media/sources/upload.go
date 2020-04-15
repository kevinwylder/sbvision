package sources

import (
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/kevinwylder/sbvision"
)

type videoSource struct {
	path  string
	title string
}

// VideoFileSource is a source from a video file
func VideoFileSource(data io.ReadCloser, title string) (sbvision.VideoSource, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, data)
	if err != nil {
		return nil, err
	}
	err = data.Close()
	if err != nil {
		return nil, err
	}
	return &videoSource{
		path:  file.Name(),
		title: title,
	}, nil
}

func (s *videoSource) GetVideo() sbvision.Video {
	return sbvision.Video{
		Title:     s.title,
		Type:      sbvision.UploadedVideo,
		SourceURL: s.path,
	}
}

func (s *videoSource) GetThumbnail() (io.ReadCloser, error) {
	command := exec.Command("ffmpeg", "-i", s.path, "-ss", "00:00:02.000", "-vframes", "1", "-f", "image2pipe", "-")
	pipe, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = command.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		command.Wait()
	}()
	return pipe, nil
}
