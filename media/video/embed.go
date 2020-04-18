package video

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/kevinwylder/sbvision"
)

// startDownload downloads and embeds a frame counter into the video
func (q *ProcessQueue) startDownload(video *sbvision.Video) *ffmpegProcess {

	dir, err := ioutil.TempDir("", "ffmpeg")
	playlist := path.Join(dir, "playlist.m3u8")
	file := path.Join(dir, "video.mp4")

	process := ffmpegProcess{
		Info:      video,
		OutputDir: dir,

		process: exec.Command("ffmpeg",
			"-i", video.SourceURL,

			"-vf", generateFfmpegFilter(16, 4, 2),
			"-an",
			"-profile:v", "main",
			"-crf", "20",
			"-g", "48", "-keyint_min", "8",
			"-sc_threshold", "0",
			"-b:v", "2500k", "-maxrate", "2675k", "-bufsize", "3750k",
			"-hls_time", "4",
			"-hls_segment_filename", fmt.Sprintf("%s/%%03d.ts", dir), playlist,

			"-vf", generateFfmpegFilter(16, 4, 2),
			"-an",
			"-profile:v", "main",
			"-crf", "20",
			"-f", "mp4",
			"-c:v", "libx264",
			file,
		),
		progress: make(chan string),
		err:      err,
	}

	if err == nil {
		go process.start(process.getDownloadProgress)
	} else {
		defer func() {
			close(process.progress)
		}()
	}

	return &process
}

func (p *ffmpegProcess) getDownloadProgress() {
	scrapeTime := regexp.MustCompile(`time=([\d:.]+)`)
	for {
		line, err := p.reader.ReadBytes('\r')
		if err != nil {
			p.err = err
			break
		}

		matches := scrapeTime.FindSubmatch(line)
		if len(matches) == 2 {
			p.progress <- string(matches[1])
		}
	}
}

func generateFfmpegFilter(bits, width, height int) string {
	var filter strings.Builder
	pow := 1
	for i := 0; i < bits; i++ {
		for j := 0; j < 2; j++ {
			filter.WriteString("drawbox=enable='eq(mod(floor(n/")
			filter.WriteString(strconv.Itoa(pow))
			filter.WriteString("),2),")
			filter.WriteString(strconv.Itoa(j))
			filter.WriteString(")':x=")
			filter.WriteString(strconv.Itoa(i * width))
			filter.WriteString(":y=0:w=")
			filter.WriteString(strconv.Itoa(width))
			filter.WriteString(":h=")
			filter.WriteString(strconv.Itoa(height))
			filter.WriteString(":color=")
			if j == 0 {
				filter.WriteString("black")
			} else {
				filter.WriteString("white")
			}
			if j == 0 || i != bits-1 {
				filter.WriteRune(',')
			}
		}
		pow *= 2
	}
	return filter.String()
}
