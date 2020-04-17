package video

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/kevinwylder/sbvision"
)

// ffmpegProcess is a struct containing information about the process embedding a frame counter in the video
type ffmpegProcess struct {
	OutputPath string
	Info       *sbvision.Video

	afterFunc func()
	process   *exec.Cmd
	reader    *bufio.Reader
	progress  chan string
	err       error
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
	p.process.Process.Kill()
	return os.Remove(p.OutputPath)
}

type hook struct {
	matcher *regexp.Regexp
	handler func([][]byte)
}

// getInfo fills in the information (resolution, duration, format, fps) from the given source
func getInfo(video *sbvision.Video) *ffmpegProcess {

	process := ffmpegProcess{
		Info:     video,
		process:  exec.Command("ffprobe", "-i", video.SourceURL),
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
