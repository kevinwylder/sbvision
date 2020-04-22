package sourceutil

import (
	"encoding/json"
	"log"
	"net/http"
	"path"

	"github.com/aws/aws-sdk-go/service/batch"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
)

// VideoRequestManager manages all the incoming http topic messages, and routes them to listeners
type VideoRequestManager struct {
	uploader     *s3manager.Uploader
	ns           *sns.SNS
	batch        *batch.Batch
	userRequests map[string]*UserRequests
	requestVideo map[string]*videoRequest
}

// NewVideoRequestManager creates a VideoRequestManager that dispatches status updates to listeners
func NewVideoRequestManager(sess *session.Session) *VideoRequestManager {
	return &VideoRequestManager{
		uploader:     s3manager.NewUploader(sess),
		ns:           sns.New(sess),
		batch:        batch.New(sess),
		userRequests: make(map[string]*UserRequests),
		requestVideo: make(map[string]*videoRequest),
	}
}

func (m *VideoRequestManager) createTopicAndSubscribe(request *videoRequest) error {
	topic, err := m.ns.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(request.ID),
	})
	if err != nil {
		return err
	}
	request.TopicARN = *topic.TopicArn
	output, err := m.ns.Subscribe(&sns.SubscribeInput{
		Protocol: aws.String("https"),
		Endpoint: aws.String("https://api.skateboardvision.net/sns/" + request.ID),
		TopicArn: topic.TopicArn,
	})
	return err
}

func (m *VideoRequestManager) handleSNSEvent(w http.ResponseWriter, r *http.Request) {
	ID := path.Base(r.URL.Path)
	request, exists := m.requestVideo[ID]
	if !exists {
		http.Error(w, "Not found", 404)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if r.Header.Get("x-amz-sns-message-type") == "SubscriptionConfirmation" {
		var body struct {
			Token string `json:"Token"` // redundantly repeating the same thing again
		}
		err := decoder.Decode(&body)
		if err != nil {
			log.Println("Error confirming SNS Subscription, bad message body - ", err.Error())
			return
		}
		_, err = m.ns.ConfirmSubscription(&sns.ConfirmSubscriptionInput{
			Token:    &body.Token,
			TopicArn: &request.TopicARN,
		})
		if err != nil {
			log.Println("Error confirming topic - ", err.Error())
		}
		return
	}

	var body struct {
		Message string `json:"Message"`
	}
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "bad request body", 400)
		return
	}
	if err := json.Unmarshal([]byte(body.Message), &request.Status); err != nil {
		log.Println("Request message was not json - ", body.Message)
		return
	}
	request.sendStatus()

	if request.Status.IsComplete {
		m.endRequest(request)
	}

}

func (m *VideoRequestManager) endRequest(request *videoRequest) {
	delete(m.requestVideo, request.ID)
}
