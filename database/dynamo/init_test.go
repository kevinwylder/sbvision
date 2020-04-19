package dynamo_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/database/dynamo"
)

func TestTableInit(t *testing.T) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		t.Fatal(err)
	}
	db, err := dynamo.FindTables(session)
	if err != nil {
		t.Fatal(err)
	}
	var user *sbvision.User
	if false {
		user = &sbvision.User{
			Email:    "wylderkevin@gmail.com",
			Username: "kwylder",
		}
		err = db.AddUser(user)
	} else {
		user, err = db.GetUser("wylderkevin@gmail.com")
	}
	if err != nil {
		t.Fatal(err)
	}
	video := &sbvision.Video{
		Duration:   "00:04:20.690",
		FPS:        30.1,
		Height:     1080,
		Width:      1920,
		Title:      "This is not a real video",
		Type:       sbvision.RedditVideo,
		UploadedAt: "2020-04-20 04:20:69",
		ShareURL:   "reddit.com/r/skateboarding",
	}
	for i := 0; i < 3; i++ {
		err = db.AddVideo(video, user)
		if err != nil {
			t.Fatal(err)
		}
	}
	videos, err := db.GetVideos(user)
	if err != nil {
		t.Fatal(err)
	}
	for i := range videos {
		fmt.Println("Deleting", videos[i])
		err := db.RemoveVideo(&videos[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log(err)
	t.Fail()
}
