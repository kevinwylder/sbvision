package encoder

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/video"
)

// videoRequest represents a request to upload a video
type videoRequest struct {
	u *UserRequests
	m *VideoRequestManager

	ID       string
	TopicARN string
	Status   video.Status

	url   string
	file  *os.File
	title string
	video sbvision.VideoType
}

// NewRequest adds a new request to the manager and user
func (u *UserRequests) NewRequest(url, title string, file *os.File) {
	data := make([]byte, 18)
	rand.Read(data)
	randID := base64.URLEncoding.EncodeToString(data)
	r := &videoRequest{
		u:  u,
		m:  u.m,
		ID: randID,
		Status: video.Status{
			RequestID: randID,
		},
		file:  file,
		title: title,
		url:   url,
	}
	u.requests = append(u.requests, r)
	go r.process()
	u.m.requestVideo[randID] = r
}

func (r *videoRequest) sendStatus() {
	for _, cb := range r.u.callbacks {
		cb(&r.Status)
	}
}

func (r *videoRequest) setStatus(status string) {
	r.Status.Message = status
	r.sendStatus()
}

func (r *videoRequest) process() {
	var err error

	defer func() {
		if err != nil {
			fmt.Println("Process function for reqest", r.ID, "exited with error", err.Error())
			r.Status.Message = err.Error()
			r.Status.IsComplete = true
			r.u.m.endRequest(r)
		}
		r.sendStatus()
	}()

	if r.file == nil {
		r.setStatus("Getting Video from Internet Source")
		err = r.downloadVideoFromInternet()
		if err != nil {
			return
		}
	}

	r.setStatus("Storing Unprocessed Video")
	err = r.uploadVideoToBucket()
	if err != nil {
		return
	}

	r.setStatus("Setting up communication channel for Video Processor")
	err = r.u.m.createTopicAndSubscribe(r)
	if err != nil {
		return
	}

	r.setStatus("Sending Request to Process Video")
	err = r.startBatchProcess()
	if err != nil {
		return
	}

	r.setStatus("Waiting for worker to run job")
}

func (r *videoRequest) downloadVideoFromInternet() error {
	var err error
	r.file, r.title, r.video, err = video.FindVideoSource(r.url)
	return err
}

func (r *videoRequest) uploadVideoToBucket() error {
	fmt.Println(r.file)
	info, err := os.Stat(r.file.Name())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(info.Size(), "bytes about to be uploaded")
	_, err = r.u.m.uploader.Upload(&s3manager.UploadInput{
		Body:   r.file,
		Key:    aws.String(r.ID),
		Bucket: aws.String(video.QueueBucket),
	})
	if err != nil {
		return err
	}
	return os.Remove(r.file.Name())
}

func (r *videoRequest) startBatchProcess() error {
	if r.TopicARN == "" || r.ID == "" || r.title == "" || r.video == -1 || r.u.user.Email == "" {
		fmt.Println(r)
		fmt.Println(r.u.user)
		return fmt.Errorf("Missing request data")
	}

	_, err := r.m.batch.SubmitJob(&batch.SubmitJobInput{
		JobDefinition: aws.String("sbgetvid"),
		JobQueue:      aws.String(video.BatchQueueName),
		JobName:       aws.String(r.ID),
		ContainerOverrides: &batch.ContainerOverrides{
			Command: []*string{
				aws.String("sbgetvid"),
				aws.String("-email=" + r.u.user.Email),
				aws.String("-request=" + r.ID),
				aws.String("-source=" + r.url),
				aws.String("-title=" + r.title),
				aws.String("-topic=" + r.TopicARN),
				aws.String(fmt.Sprintf("-type=%d", r.video)),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
