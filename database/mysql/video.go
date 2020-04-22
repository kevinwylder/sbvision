package mysqldb

import (
	"time"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddVideo() (err error) {
	sb.addVideo, err = sb.db.Prepare(`
INSERT INTO videos (title, width, height, fps, duration, type, uploaded_by, upload_time, share_url, source_url)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	return
}

// AddVideo adds the video to the given user account
func (sb *SBDatabase) AddVideo(video *sbvision.Video, user *sbvision.User) error {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	if video.Type == sbvision.YoutubeVideo {
		video.SourceURL = video.ShareURL
	}
	result, err := sb.addVideo.Exec(video.Title, video.Width, video.Height, video.FPS, video.Duration, video.Type, user.ID, now, video.ShareURL, video.SourceURL)
	if err != nil {
		return err
	}
	video.UploadedAt = now
	video.UploadedBy = user.Username
	video.ID, err = result.LastInsertId()
	return err
}
func (sb *SBDatabase) prepareGetVideoByID() (err error) {
	sb.getVideoByID, err = sb.db.Prepare(`
SELECT videos.id, title, width, height, fps, duration, type, upload_time, users.username, share_url, source_url
FROM videos 
INNER JOIN users ON users.id = videos.uploaded_by
WHERE videos.id = ?;
	`)
	return
}

// GetVideoByID gets the video and the uploader id
func (sb *SBDatabase) GetVideoByID(id int64) (*sbvision.Video, error) {
	var video sbvision.Video
	result := sb.getVideoByID.QueryRow(id)
	err := result.Scan(&video.ID, &video.Title, &video.Width, &video.Height, &video.FPS, &video.Duration, &video.Type, &video.UploadedAt, &video.UploadedBy, &video.ShareURL, &video.SourceURL)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (sb *SBDatabase) prepareGetVideos() (err error) {
	sb.getVideos, err = sb.db.Prepare(`
SELECT id, title, width, height, fps, duration, type, upload_time, share_url, source_url
FROM videos 
WHERE uploaded_by = ?
ORDER BY upload_time DESC;
	`)
	return
}

// GetVideos gets all the videos this user has uploaded
func (sb *SBDatabase) GetVideos(user *sbvision.User) ([]sbvision.Video, error) {
	results, err := sb.getVideos.Query(user.ID)
	if err != nil {
		return nil, err
	}
	var videos []sbvision.Video
	for results.Next() {
		var video sbvision.Video
		err = results.Scan(&video.ID, &video.Title, &video.Width, &video.Height, &video.FPS, &video.Duration, &video.Type, &video.UploadedAt, &video.ShareURL, &video.SourceURL)
		video.UploadedBy = user.Username
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (sb *SBDatabase) prepareRemoveVideo() (err error) {
	sb.removeVideo, err = sb.db.Prepare(`
DELETE FROM videos WHERE id = ?
	`)
	return
}

// RemoveVideo removes this video from the account
func (sb *SBDatabase) RemoveVideo(video *sbvision.Video) error {
	_, err := sb.removeVideo.Exec(video.ID)
	return err
}