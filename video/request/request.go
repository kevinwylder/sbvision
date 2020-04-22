package sourceutil

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/batch"

	"github.com/aws/aws-sdk-go/aws"
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
func (u *UserRequests) NewRequest(url, title string, file *os.File) error {
	data := make([]byte, 18)
	_, err := rand.Read(data)
	if err != nil {
		return err
	}
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
	go r.process()
	u.m.requestVideo[randID] = r
	return nil
}

func (r *videoRequest) sendStatus() {
	for _, cb := range r.u.callbacks {
		cb(&r.Status)
	}
}

func (r *videoRequest) process() {
	var err error

	defer func() {
		if err != nil {
			r.Status.Message = err.Error()
			r.Status.IsComplete = true
			r.u.m.endRequest(r)
		}
		r.sendStatus()
	}()

	if r.file == nil {
		err = r.downloadVideoFromInternet()
		if err != nil {
			return
		}
	}

	err = r.uploadVideoToBucket()
	if err != nil {
		return
	}

	err = r.u.m.createTopicAndSubscribe(r)
	if err != nil {
		return
	}

	err = r.startBatchProcess()
	if err != nil {
		return
	}
}

func (r *videoRequest) downloadVideoFromInternet() error {
	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r.file, err = ioutil.TempFile("", "")
	if err != nil {
		return err
	}

	_, err = io.Copy(r.file, resp.Body)
	return err
}

func (r *videoRequest) uploadVideoToBucket() error {
	_, err := r.u.m.uploader.Upload(&s3manager.UploadInput{
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
	r.u.user.Email
	r.m.batch.SubmitJob(&batch.SubmitJobInput{})

	return nil
}
