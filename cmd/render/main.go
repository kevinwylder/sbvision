package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
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
	sb, err := skateboard.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "frames")

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-s
		os.RemoveAll(dir)
		sb.Destroy()
	}()

	fmt.Println("rendering at", dir)
	for i, arg := range os.Args[1:] {
		var rotation [4]float64
		err = json.Unmarshal([]byte(arg), &rotation)
		if err != nil {
			fmt.Printf(arg, "(position %d) is not a JSON array of 4 values\n", i)
			break
		}

		image := image.NewRGBA(image.Rect(0, 0, 500, 500))
		data := sb.Render(rotation)
		for x := 0; x < 500; x++ {
			for y := 0; y < 500; y++ {
				i := 4 * (((499 - y) * 500) + x)
				if data[i+3] != 0 {
					image.SetRGBA(x, y, color.RGBA{
						R: data[i],
						G: data[i+1],
						B: data[i+2],
						A: 255,
					})
				} else {
					image.Set(x, y, color.White)
				}
			}
		}

		file, err := os.Create(path.Join(dir, fmt.Sprintf("frame%03d.bmp", i)))
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

	fmt.Println("Encoding video")
	os.Remove("out.mp4")
	data, err := exec.Command(
		"ffmpeg",
		"-framerate", "1",
		"-i", path.Join(dir, "frame%03d.bmp"),
		"-c:v", "libx264",
		"out.mp4",
	).CombinedOutput()
	if err != nil {
		fmt.Println(string(data))
	}
	os.RemoveAll(dir)
	sb.Destroy()
}
