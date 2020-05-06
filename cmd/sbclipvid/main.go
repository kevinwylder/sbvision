package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database/dynamo"
	"github.com/kevinwylder/sbvision/video/skateboard"
)

type runtime struct {
	clip       *sbvision.Clip
	video      *sbvision.Video
	workdir    string
	output     string
	skateboard *skateboard.Renderer
	ddb        *dynamo.SBDatabase
}

func main() {

	var (
		clipID = flag.String("clip", "", "the id of the clip to process")
		output = flag.String("out", "", "the ouput file to write")
	)
	flag.Parse()

	if *clipID == "" || *output == "" {
		fmt.Println("All args are required!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ddb, err := dynamo.FindTables(sess)
	if err != nil {
		log.Fatal(err)
	}

	clip, err := ddb.GetClipByID(*clipID)
	if err != nil {
		log.Fatal(err)
	}

	video, err := ddb.GetVideoByID(clip.VideoID)
	if err != nil {
		log.Fatal(err)
	}

	sb, err := skateboard.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	runtime := runtime{
		clip:       clip,
		video:      video,
		output:     *output,
		ddb:        ddb,
		workdir:    dir,
		skateboard: sb,
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

func (rt *runtime) Begin() error {
	err := rt.getVideoFrames()
	if err != nil {
		return err
	}
	err = rt.transformFrames()
	if err != nil {
		return err
	}
	err = rt.encodeVideoFrames()
	return err
}

func (rt *runtime) Cleanup() {
	os.RemoveAll(rt.workdir)
}
