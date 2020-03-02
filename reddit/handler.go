package reddit

import (
	"net/http"

	"github.com/kevinwylder/sbvision"
)

type RedditHandler struct {
	youtube sbvision.VideoHandler
	client  http.Client
}

func NewRedditHandler(youtube sbvision.VideoHandler) *RedditHandler {
	handler := RedditHandler{
		youtube: youtube,
	}
}
