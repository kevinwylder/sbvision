package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kevinwylder/sbvision/cdn"
)

func (rt *runtime) getVideoInformation() error {
	process := rt.getInfo()
	for range <-process.Progress() {
	}
	return process.Error()
}

func (rt *runtime) getThumbnail() error {
	file, err := os.Create(path.Join(rt.tmpdir, "thumbnail.jpg"))
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
	}()
	command := exec.Command("ffmpeg", "-i", rt.file.Name(), "-ss", "00:00:02.000", "-vframes", "1", "-f", "image2pipe", "-")
	fmt.Println(strings.Join(command.Args, " "))
	pipe, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	err = command.Start()
	if err != nil {
		return err
	}
	go func() {
		command.Wait()
	}()
	_, err = io.Copy(file, pipe)
	return err
}

func (rt *runtime) processVideo() error {
	process := rt.embedFrameCounterAndDownsample()
	for status := range process.Progress() {
		rt.setStatus("Encoding video. Progress " + status + " of " + rt.video.Duration)
	}
	if err := process.Error(); err != nil {
		return err
	}
	rt.tmpdir = process.OutputDir
	return nil
}

func (rt *runtime) addToDatabase() error {
	return rt.db.AddVideo(rt.video, rt.user)
}

func (rt *runtime) uploadVideo() error {
	dst := cdn.VideoDirectory(rt.video.ID)
	data, err := os.Create(path.Join(rt.tmpdir, "data.json"))
	if err != nil {
		return err
	}
	err = json.NewEncoder(data).Encode(rt.video)
	if err != nil {
		return err
	}
	data.Close()
	err = rt.cdn.AddDir(rt.tmpdir, dst)
	if err != nil {
		return err
	}
	defer os.RemoveAll(rt.tmpdir)
	return rt.cdn.Invalidate(path.Join(dst, "*"))
}
