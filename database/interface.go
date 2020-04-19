package database

import "github.com/kevinwylder/sbvision"

// SBDatabase is an interface wrapping a data provider
type SBDatabase interface {
	AddUser(user *sbvision.User) error
	GetUser(email string) (*sbvision.User, error)

	AddVideo(video *sbvision.Video, user *sbvision.User) error
	GetVideoByID(id int64) (*sbvision.Video, error)
	GetVideos(user *sbvision.User) ([]sbvision.Video, error)
	RemoveVideo(video *sbvision.Video) error
}
