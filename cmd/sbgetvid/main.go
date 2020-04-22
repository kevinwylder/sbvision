package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kevinwylder/sbvision"

	"github.com/kevinwylder/sbvision/cdn"
	"github.com/kevinwylder/sbvision/database/dynamo"
	"github.com/kevinwylder/sbvision/video"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
)

type runtime struct {
	user  *sbvision.User
	video *sbvision.Video
	db    *dynamo.SBDatabase

	request string
	topic   *string
	ns      *sns.SNS
	status  video.Status

	file   *os.File
	tmpdir string

	cdn *cdn.Uploader
}

// sbgetvid gets a video from the internet and encodes it to mobile friendly formats
// there is at least 1 HLS format, and an mp4 complete file
// the video will have a frame counter embedded in the top 2 pixels of the video
func main() {
	var (
		snsTopic  = flag.String("topic", "", "pass the sns arn to publish to")
		requestID = flag.String("request", "", "pass the path in the s3 bucket to pull from")
		sourceURL = flag.String("source", "", "url that this came from (reddit/youtube)")
		title     = flag.String("title", "", "pass the video title")
		videoType = flag.Int("type", -1, "pass the sbvision.VideoType type ")
		userEmail = flag.String("email", "", "Pass the users email address")
	)
	flag.Parse()

	if *snsTopic == "" || *requestID == "" || *title == "" || *videoType == -1 || *userEmail == "" {
		flag.PrintDefaults()
		log.Fatal()
	}

	dir, err := ioutil.TempDir("", "ffmpeg")
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

	user, err := ddb.GetUser(*userEmail)
	if err != nil {
		log.Fatal(err)
	}

	v := sbvision.Video{
		Title:         *title,
		Type:          sbvision.VideoType(*videoType),
		SourceURL:     *sourceURL,
		UploaderEmail: user.Email,
	}

	rt := &runtime{
		db:      ddb,
		user:    user,
		video:   &v,
		request: *requestID,
		status: video.Status{
			RequestID: *requestID,
			Video:     &v,
		},
		ns:     sns.New(sess),
		topic:  snsTopic,
		cdn:    cdn.NewUploader(sess),
		tmpdir: dir,
	}

	rt.downloadVideo(s3manager.NewDownloader(sess))
	rt.embed()
	rt.cleanup(s3.New(sess))
}

func (rt *runtime) downloadVideo(manager *s3manager.Downloader) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	rt.file = f
	_, err = manager.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(video.QueueBucket),
		Key:    aws.String(rt.request),
	})
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	stat, err := os.Stat(f.Name())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded", stat.Size(), "bytes to", f.Name())
}

func (rt *runtime) embed() {

	defer func() {
		rt.finish()
	}()

	rt.setStatus("Getting video information")
	if err := rt.getVideoInformation(); err != nil {
		rt.setStatus("Error getting video info: " + err.Error())
		return
	}

	rt.setStatus("Creating Thumbnail")
	if err := rt.getThumbnail(); err != nil {
		rt.setStatus("Failed to get thumbnail - " + err.Error())
		return
	}

	rt.setStatus("Processing Video")
	if err := rt.processVideo(); err != nil {
		rt.setStatus("Error processing video - " + err.Error())
	}

	rt.setStatus("Adding to the database")
	if err := rt.addToDatabase(); err != nil {
		rt.setStatus("Failed to add to the database - " + err.Error())
		return
	}

	rt.setStatus("Uploading video to skateboardvision.net")
	if err := rt.uploadVideo(); err != nil {
		rt.setStatus("Error uploading video - " + err.Error())
	}

	rt.setStatus("Complete")
	rt.status.WasSuccess = true

}

func (rt *runtime) cleanup(s3sess *s3.S3) {
	//os.Remove(rt.file.Name())
	//os.RemoveAll(rt.tmpdir)
	s3sess.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(video.QueueBucket),
		Key:    aws.String(rt.request),
	})
	rt.ns.DeleteTopic(&sns.DeleteTopicInput{
		TopicArn: rt.topic,
	})
}
