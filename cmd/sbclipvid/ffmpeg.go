package main

import (
	"fmt"
	"image/png"
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
		path.Join(rt.output, "frame%03d.png"),
	).Run()
}

func (rt *runtime) computeInterpolation() error {
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
		err := rt.createSubFrames(frame, maxDimension)
		if err != nil {
			return err
		}
	}

	return nil
}

const subframeCount int64 = 16

func (rt *runtime) createSubFrames(frame int64, maxDimension int64) error {
	in, err := os.Open(path.Join(rt.output, fmt.Sprintf("frame%03d.png", frame+1)))
	if err != nil {
		return err
	}
	defer in.Close()
	inBmp, err := png.Decode(in)
	if err != nil {
		return err
	}
	bound := rt.clip.Bounds[rt.clip.Start+frame]
	var subframe int64
	for subframe = 0; subframe < subframeCount; subframe++ {
		time := float64(frame) + float64(subframe)/float64(subframeCount)
		percent := time / float64(rt.clip.End-rt.clip.Start+1)
		image := rt.drawImageV1(inBmp, bound, maxDimension, rt.function.At(percent*rt.function.Duration()))
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
	return exec.Command(
		"ffmpeg",
		"-framerate", fmt.Sprint(60),
		"-i", path.Join(rt.workdir, "output_%03d.bmp"),
		"-c:v", "libx264",
		path.Join(rt.output, "clip.mp4"),
	).Run()
}
