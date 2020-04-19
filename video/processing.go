package video

import "github.com/kevinwylder/sbvision"

// QueueBucket is the name of the bucket to put source videos
const QueueBucket string = "skateboardvision-videos"

// BatchQueueName is the name of the aws batch processing queue to target
const BatchQueueName string = "video-processing-queue"

// Status is a status update from a video request
type Status struct {
	RequestID  string          `json:"requestid"`
	Video      *sbvision.Video `json:"info"`
	Message    string          `json:"message"`
	IsComplete bool            `json:"is_complete"`
	WasSuccess bool            `json:"was_success"`
}
