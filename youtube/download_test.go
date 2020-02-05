package youtube_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/youtube"
)

type mockResults struct {
	addedThumbnailImage bool
	uploadedThumbnail   bool
	addedVideo          bool
	addedYoutubeVideo   bool
	test                *testing.T
}

func TestYoutubeDownload(t *testing.T) {
	mock := &mockResults{
		test: t,
	}
	handler := youtube.NewYoutubeHandler(mock, mock)
	_, err := handler.HandleDiscover(&sbvision.VideoDiscoverRequest{
		Session: &sbvision.Session{ID: 1, IP: "192.168.toaster", Time: time.Now().Unix()},
		Type:    1,
		URL:     "https://www.youtube.com/watch?v=ardIYGi_Ras",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !mock.uploadedThumbnail || !mock.addedThumbnailImage || !mock.addedVideo || !mock.addedYoutubeVideo {
		t.Fail()
	}
}

func (r *mockResults) PutImage(data io.Reader, key sbvision.Image) error {
	if r.uploadedThumbnail {
		r.test.Fail()
	}
	r.uploadedThumbnail = true
	fmt.Println("upload image")
	return nil
}

func (r *mockResults) AddImage(image sbvision.Image, session *sbvision.Session) error {
	if !r.uploadedThumbnail {
		r.test.Fail()
	}
	if r.addedThumbnailImage {
		r.test.Fail()
	}
	fmt.Println("add image")
	r.addedThumbnailImage = true
	return nil
}

func (r *mockResults) AddVideo(video *sbvision.Video) error {
	fmt.Println("add video")
	if !r.uploadedThumbnail || !r.addedThumbnailImage {
		r.test.Fail()
	}
	if r.addedVideo {
		r.test.Fail()
	}
	r.addedVideo = true
	return nil
}

func (r *mockResults) AddYoutubeRecord(video *sbvision.YoutubeVideoInfo) error {
	if !r.addedVideo {
		r.test.Fail()
	}
	if r.addedYoutubeVideo {
		r.test.Fail()
	}
	fmt.Println("add youtube video")
	r.addedYoutubeVideo = true
	return nil
}

func (r *mockResults) GetYoutubeRecord(videoID int64) (*sbvision.YoutubeVideoInfo, error) {
	r.test.Fail()
	return nil, nil
}
