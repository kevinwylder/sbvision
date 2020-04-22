package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

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
		rt.setStatus(status)
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
	files, err := ioutil.ReadDir(rt.tmpdir)
	if err != nil {
		return err
	}
	for _, file := range files {
		f, err := os.Open(path.Join(rt.tmpdir, file.Name()))
		if err != nil {
			return err
		}
		err = rt.cdn.Add(f, path.Join(cdn.VideoDirectory(rt.video.ID), file.Name()))
		if err != nil {
			return err
		}
		f.Close()
		os.Remove(f.Name())
	}
	return rt.cdn.Invalidate(path.Join(cdn.VideoDirectory(rt.video.ID), "*"))
}
