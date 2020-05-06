package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"golang.org/x/image/bmp"
)

func (rt *runtime) getVideoFrames() error {
	return exec.Command(
		"ffmpeg",
		"-i", fmt.Sprintf("https://skateboardvision.net/video/%s/video.mp4", rt.clip.VideoID),
		"-vf", fmt.Sprintf("trim=start_frame=%d:end_frame=%d", rt.clip.Start, rt.clip.End+1),
		path.Join(rt.workdir, "input_%03d.bmp"),
	).Run()
}

func (rt *runtime) transformFrames() error {
	var frame int64
	var maxDimension int64
	for _, bound := range rt.clip.Bounds {
		if bound.Width > maxDimension {
			maxDimension = bound.Width
		}
		if bound.Height > maxDimension {
			maxDimension = bound.Height
		}
	}
	for frame = 0; frame <= rt.clip.End-rt.clip.Start; frame++ {
		err := rt.createFrames(frame, maxDimension)
		if err != nil {
			return err
		}
	}
	return nil
}

const subframeCount int64 = 16

func (rt *runtime) createFrames(frame int64, maxDimension int64) error {
	in, err := os.Open(path.Join(rt.workdir, fmt.Sprintf("input_%03d.bmp", frame+1)))
	if err != nil {
		return err
	}
	defer in.Close()
	inBmp, err := bmp.Decode(in)
	if err != nil {
		return err
	}
	bound := rt.clip.Bounds[rt.clip.Start+frame]
	var subframe int64
	for subframe = 0; subframe < subframeCount; subframe++ {
		time := float64(frame+rt.clip.Start) + float64(subframe)/float64(subframeCount)
		image := rt.drawImageV1(inBmp, bound, maxDimension, interpolateRotation(rt.clip, time))
		if err != nil {
			return err
		}
		out, err := os.Create(path.Join(rt.workdir, fmt.Sprintf("output_%03d.bmp", frame*subframeCount+subframe)))
		if err != nil {
			return err
		}
		err = bmp.Encode(out, image)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

func (rt *runtime) encodeVideoFrames() error {
	os.Remove(rt.output)
	return exec.Command(
		"ffmpeg",
		"-framerate", fmt.Sprint(rt.video.FPS*2),
		"-i", path.Join(rt.workdir, "output_%03d.bmp"),
		"-c:v", "libx264",
		rt.output,
	).Run()
}
