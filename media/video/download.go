package video

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/kevinwylder/sbvision"
)

// FfmpegProcess is a struct containing information about the process embedding a frame counter in the video
type FfmpegProcess struct {
	OutputPath string
	Info       sbvision.Video

	process  *exec.Cmd
	reader   *bufio.Reader
	progress chan string
	err      error
}

// StartDownload downloads and embeds a frame counter into the video
func StartDownload(source string) (*FfmpegProcess, error) {
	tmp, err := ioutil.TempFile("sbvisionvideo", source)
	if err != nil {
		return nil, err
	}
	err = tmp.Close()
	if err != nil {
		return nil, err
	}

	process := FfmpegProcess{
		OutputPath: tmp.Name(),
		process:    exec.Command("ffmpeg", "-i", source, "-vf", generateFfmpegFilter(16, 4, 2), "-y", tmp.Name()),
		progress:   make(chan string),
	}

	go process.start()

	return &process, nil
}

// Error returns the error associated with the process if there was any
func (p *FfmpegProcess) Error() error {
	return p.err
}

// Progress returns the progress channel of the process
// when the process exits, this channel is closed
func (p *FfmpegProcess) Progress() <-chan string {
	return p.progress
}

// Cancel stops the process of embedding and removes the video
func (p *FfmpegProcess) Cancel() error {
	p.process.Process.Kill()
	return os.Remove(p.OutputPath)
}

func (p *FfmpegProcess) start() {
	p.err = p.process.Start()
	if p.err != nil {
		close(p.progress)
		return
	}

	reader, err := p.process.StderrPipe()
	if err != nil {
		p.err = err
		close(p.progress)
		return
	}

	p.reader = bufio.NewReaderSize(reader, 1024)

	// read the stream info
	err = p.readInfo()
	if err != nil {
		p.err = fmt.Errorf("%s\ndid not finish parsing video info", err.Error())
		p.process.Wait()
		close(p.progress)
		return
	}

	// push the progress to the stream
	scrapeTime := regexp.MustCompile(`time=([\d:.]+)`)
	for {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			reader.Close()
			break
		}

		matches := scrapeTime.FindSubmatch(line)
		if len(matches) == 2 {
			p.progress <- string(matches[1])
		}
	}

	p.err = p.process.Wait()
	close(p.progress)

}
