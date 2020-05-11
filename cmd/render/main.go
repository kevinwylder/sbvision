package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"

	"golang.org/x/image/bmp"

	"github.com/kevinwylder/sbvision/video/skateboard"
)

func main() {

	var (
		video = flag.Bool("video", false, "if true, creates a video of the frames")
		fps   = flag.Float64("fps", 30.0, "if video is set, use this framerate")
		name  = flag.String("name", "out", "the name of the asset(s) to generate")
	)
	flag.Parse()

	sb, err := skateboard.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	dir := "."
	if *video {
		dir, err = ioutil.TempDir("", "")
		if err != nil {
			sb.Destroy()
			log.Fatal(err)
		}
	}

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-s
		sb.Destroy()
		if *video {
			os.RemoveAll(dir)
		}
	}()

	for i, arg := range flag.Args() {
		var rotation [4]float64
		err = json.Unmarshal([]byte(arg), &rotation)
		if err != nil {
			fmt.Printf(arg, "(position %d) is not a JSON array of 4 values\n", i)
			break
		}

		image := sb.Render(rotation)

		fileName := path.Join(dir, fmt.Sprintf("%s%03d.bmp", *name, i))
		fmt.Println(fileName)
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Could not create file")
			break
		}

		err = bmp.Encode(file, image)
		if err != nil {
			fmt.Println("Error encoding bitmap")
			break
		}

	}

	sb.Destroy()

	if *video {
		videoName := *name + ".mp4"
		fmt.Println("Encoding ", videoName)
		os.Remove(videoName)
		data, err := exec.Command(
			"ffmpeg",
			"-framerate", fmt.Sprint(*fps),
			"-i", path.Join(dir, *name+"%03d.bmp"),
			"-c:v", "libx264",
			videoName,
		).CombinedOutput()
		if err != nil {
			fmt.Println(string(data))
		}

		os.RemoveAll(dir)
	}
}
