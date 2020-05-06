package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database/dynamo"
)

type runtime struct {
	clip       *sbvision.Clip
	video      *sbvision.Video
	workdir    string
	output     string
	skateboard *skateboard
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

	sb, err := newSkateboard()
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
	defer os.RemoveAll(rt.workdir)
	files, err := ioutil.ReadDir(rt.workdir)
	if err != nil {
		return
	}
	dst := path.Dir(rt.output)
	for _, f := range files {
		if strings.Contains(f.Name(), "output") {
			in, _ := os.Open(path.Join(rt.workdir, f.Name()))
			out, _ := os.OpenFile(path.Join(dst, f.Name()), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
			io.Copy(out, in)
			in.Close()
			out.Close()
		}
	}
}
