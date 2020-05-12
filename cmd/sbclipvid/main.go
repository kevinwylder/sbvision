package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/cdn"
	"github.com/kevinwylder/sbvision/video/interpolate"
	"github.com/kevinwylder/sbvision/video/skateboard"
)

type runtime struct {
	clip       sbvision.Clip
	workdir    string
	output     string
	skateboard *skateboard.Renderer
	function   interpolate.QuaternionFunction
	cdn        *cdn.Uploader
}

func main() {
	var runtime runtime
	var err error

	if len(os.Args) < 2 {
		fmt.Println(`USAGE: sbclipvid [clip data json base64 urlencoded]`)
		os.Exit(1)
	}

	data, err := base64.URLEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &runtime.clip)
	if err != nil {
		log.Fatal(err)
	}

	runtime.function = interpolate.LowPassFilter(&runtime.clip)

	runtime.skateboard, err = skateboard.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		log.Fatal(err)
	}
	runtime.cdn = cdn.NewUploader(sess)

	runtime.workdir, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	runtime.output = path.Join(runtime.workdir, "output")
	err = os.Mkdir(runtime.output, 0777)
	if err != nil {
		log.Fatal(err)
	}

	err = runtime.Begin()
	if err != nil {
		fmt.Println(err)
	}
	runtime.Cleanup()
	if err != nil {
		os.Exit(1)
	}
}

func (rt *runtime) writeDataFile() error {
	file, err := os.Create(path.Join(rt.output, "data.json"))
	if err != nil {
		return err
	}
	var frames []sbvision.Frame
	for i := rt.clip.Start; i <= rt.clip.End; i++ {
		frames = append(frames, sbvision.Frame{
			Bound:    rt.clip.Bounds[i],
			Rotation: rt.clip.Rotations[i],
			Image:    fmt.Sprintf("/clip/%s/frame%03d.png", rt.clip.ID, i-rt.clip.Start+1),
		})
	}
	err = json.NewEncoder(file).Encode(frames)
	if err != nil {
		return err
	}
	return file.Close()
}

func (rt *runtime) Begin() error {
	err := rt.getVideoFrames()
	if err != nil {
		return err
	}
	err = rt.computeInterpolation()
	if err != nil {
		return err
	}
	err = rt.encodeVideoFrames()
	if err != nil {
		return err
	}
	err = rt.writeDataFile()
	if err != nil {
		return err
	}
	dst := cdn.ClipDirectory(rt.clip.ID)
	err = rt.cdn.AddDir(rt.output, dst)
	if err != nil {
		return err
	}
	err = rt.cdn.Invalidate(path.Join(dst, "*"))
	if err != nil {
		return err
	}
	return err
}

func (rt *runtime) Cleanup() {
	os.RemoveAll(rt.workdir)
}
