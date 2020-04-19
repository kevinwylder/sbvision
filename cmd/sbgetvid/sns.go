package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/sns"
)

func (rt *runtime) setStatus(status string) {
	log.Println(status)

	rt.status.Message = status
	data, _ := json.Marshal(&rt.status)
	message := string(data)
	_, err := rt.ns.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: rt.topic,
	})
	if err != nil {
		log.Println("PUBLISH ERROR - ", err)
	}
}

func (rt *runtime) finish() {
	rt.status.IsComplete = true
	rt.setStatus(rt.status.Message)
}
