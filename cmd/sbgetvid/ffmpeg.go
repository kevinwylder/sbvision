package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/kevinwylder/sbvision"
)

// ffmpegProcess is a struct containing information about the process embedding a frame counter in the video
type ffmpegProcess struct {
	Info      *sbvision.Video
	OutputDir string

	process  *exec.Cmd
	reader   *bufio.Reader
	progress chan string
	err      error
}

func (p *ffmpegProcess) start(parser func()) {
	reader, err := p.process.StderrPipe()
	if err != nil {
		p.err = err
		return
	}

	p.reader = bufio.NewReaderSize(reader, 1024*10)

	p.err = p.process.Start()
	if p.err != nil {
		close(p.progress)
		return
	}

	go parser()

	err = p.process.Wait()
	if err != nil {
		p.err = err
	}
	close(p.progress)

}

// Error returns the error associated with the process if there was any
func (p *ffmpegProcess) Error() error {
	return p.err
}

// Progress returns the progress channel of the process
// when the process exits, this channel is closed
func (p *ffmpegProcess) Progress() <-chan string {
	return p.progress
}

// Cancel stops the process of embedding and removes the video
func (p *ffmpegProcess) Cancel() error {
	err := p.process.Process.Kill()
	if p.OutputDir != "" {
		err = os.RemoveAll(p.OutputDir)
	}
	return err
}

type hook struct {
	matcher *regexp.Regexp
	handler func([][]byte)
}

// getInfo fills in the information (resolution, duration, format, fps) from the given source
func (rt *runtime) getInfo() *ffmpegProcess {

	process := ffmpegProcess{
		Info:     rt.video,
		process:  exec.Command("ffprobe", "-i", rt.file.Name()),
		progress: make(chan string),
	}

	go process.start(process.readInfo)

	return &process
}

func (p *ffmpegProcess) readInfo() {

	hooks := []hook{
		hook{
			matcher: regexp.MustCompile(`Duration: (\d+:\d+:\d+\.\d+),`),
			handler: func(data [][]byte) {
				p.Info.Duration = string(data[1])
			},
		},
		hook{
			matcher: regexp.MustCompile(`Stream #\d+:.*: Video:.*, (\d+)x(\d+).*, (\d+.?\d*) fps`),
			handler: func(data [][]byte) {
				p.Info.Width, p.err = strconv.ParseInt(string(data[1]), 10, 64)
				if p.err != nil {
					return
				}
				p.Info.Height, p.err = strconv.ParseInt(string(data[2]), 10, 64)
				if p.err != nil {
					return
				}
				p.Info.FPS, p.err = strconv.ParseFloat(string(data[3]), 64)
			},
		},
	}

	var line []byte
	i := 0
	for i < len(hooks) {
		line, p.err = p.reader.ReadBytes('\n')
		if p.err != nil {
			if i != len(hooks) {
				p.err = fmt.Errorf("Did not capture all info (%d of %d)", i, len(hooks))
			}
			return
		}

		matches := hooks[i].matcher.FindSubmatch(line)
		if len(matches) != 0 {
			hooks[i].handler(matches)
			if p.err != nil {
				return
			}
			i++
		}
	}
}

// embedFrameCounterAndDownsample downloads and embeds a frame counter into the video
func (rt *runtime) embedFrameCounterAndDownsample() *ffmpegProcess {
	playlist := path.Join(rt.tmpdir, "playlist.m3u8")
	file := path.Join(rt.tmpdir, "video.mp4")

	process := ffmpegProcess{
		Info:      rt.video,
		OutputDir: rt.tmpdir,

		process: exec.Command("ffmpeg",
			"-i", rt.file.Name(),

			"-vf", generateFfmpegFilter(16, 4, 2),
			"-an",
			"-profile:v", "main",
			"-crf", "20",
			"-g", "48", "-keyint_min", "8",
			"-sc_threshold", "0",
			"-b:v", "2500k", "-maxrate", "2675k", "-bufsize", "3750k",
			"-hls_time", "4",
			"-hls_segment_filename", fmt.Sprintf("%s/%%03d.ts", rt.tmpdir), playlist,

			"-vf", generateFfmpegFilter(16, 4, 2),
			"-an",
			"-profile:v", "main",
			"-crf", "20",
			"-f", "mp4",
			"-c:v", "libx264",
			file,
		),
		progress: make(chan string),
	}

	go process.start(process.getDownloadProgress)

	return &process
}

func (p *ffmpegProcess) getDownloadProgress() {
	scrapeTime := regexp.MustCompile(`time=([\d:.]+)`)
	for {
		line, err := p.reader.ReadBytes('\r')
		if err != nil {
			if err != io.EOF {
				p.err = err
			}
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
