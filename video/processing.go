package video

import "github.com/kevinwylder/sbvision"

// QueueBucket is the name of the bucket to put source videos
const QueueBucket string = "skateboardvision-videos"

// Status is a status update from a video request
type Status struct {
	RequestID  string          `json:"requestid"`
	Video      *sbvision.Video `json:"info"`
	Message    string          `json:"message"`
	IsComplete bool            `json:"is_complete"`
	WasSuccess bool            `json:"was_success"`
}
