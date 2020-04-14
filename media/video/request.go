package video

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/kevinwylder/sbvision"
)

// ProcessRequest is a request to process the given video source
type ProcessRequest struct {
	user    *sbvision.User
	url     string
	source  sbvision.VideoSource
	onEvent chan struct{}

	Info       *sbvision.Video `json:"info"`
	ID         string          `json:"id"`
	Status     string          `json:"status"`
	IsComplete bool            `json:"complete"`
	WasSuccess bool            `json:"success"`
}

// Enqueue adds this source to the list of videos to process
func (q *ProcessQueue) Enqueue(user *sbvision.User, url string) (*ProcessRequest, error) {
	if _, exists := q.Find(user); exists {
		return nil, fmt.Errorf("Cannot do multiple requests at the same time by the same user")
	}
	data := make([]byte, 10)
	rand.Reader.Read(data)
	request := &ProcessRequest{
		user:    user,
		url:     url,
		ID:      base64.URLEncoding.EncodeToString(data),
		onEvent: make(chan struct{}),
	}

	select {
	case q.queue <- request:
		q.requests[user.ID] = request
		request.setStatus("Added to queue")
		return request, nil
	default:
		return nil, fmt.Errorf("Queue is full, please come back later")
	}
}

// Find looks up the request for the given email
func (q *ProcessQueue) Find(user *sbvision.User) (*ProcessRequest, bool) {
	request, exists := q.requests[user.ID]
	if !exists || request.IsComplete {
		return nil, false
	}
	return request, true
}

func (r *ProcessRequest) finish(q *ProcessQueue) {
	r.IsComplete = true
	tmp := r.onEvent
	r.onEvent = make(chan struct{})
	close(tmp)
	defer func() {
		time.Sleep(time.Minute)
		if request, exists := q.requests[r.user.ID]; exists && request.ID == r.ID {
			delete(q.requests, r.user.ID)
		}
	}()
}

func (r *ProcessRequest) setStatus(status string) {
	r.Status = status
	tmp := r.onEvent
	r.onEvent = make(chan struct{})
	close(tmp)
}

// Subscribe creates a channel that sends on events.
// Make sure to call the cleanup func when done.
func (r *ProcessRequest) Subscribe() (<-chan struct{}, func()) {
	events := make(chan struct{})
	cancel := make(chan struct{})
	go func() {
		running := true
		for running {
			select {
			case <-r.onEvent: // always a close. Next time around it will be a new chan. Not guaranteed to capture every event.
				events <- struct{}{}
			case <-cancel:
				running = false
			}
		}
		close(events)
	}()
	return events, func() {
		close(cancel)
	}
}
