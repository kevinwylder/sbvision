package database

import (
	"database/sql"

	"github.com/kevinwylder/sbvision"
)

// GetVideos is a paged listing of videos
func (db *SBDatabase) GetVideos(offset, count int) ([]sbvision.Video, error) {
	return db.getVideos("", `
ORDER BY videos.discovery_time DESC
LIMIT ? OFFSET ?
	`, offset, count)
}

func (db *SBDatabase) getVideos(where string, order string, args ...interface{}) ([]sbvision.Video, error) {
	results, err := db.Query(`
SELECT	
	videos.id,
	videos.title,
	images.key,
	videos.type,
	videos.duration,
	videos.fps,
	COUNT(*),
	MAX(clips.id)
FROM videos
INNER JOIN images 
		ON images.id = videos.thumbnail_id
LEFT JOIN frames
		ON frames.video_id = videos.id
LEFT JOIN clips
		ON clips.frame_id = frames.id
`+where+`
GROUP BY videos.id
`+order, args...)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	var videos []sbvision.Video
	for results.Next() {
		var videoID, videoType, videoDuration, clipCount int64
		var fps float64
		var s3Key, title string
		var clipFound sql.NullInt64

		err := results.Scan(
			&videoID,
			&title,
			&s3Key,
			&videoType,
			&videoDuration,
			&fps,
			&clipCount,
			&clipFound,
		)
		if err != nil {
			return nil, err
		}

		if !clipFound.Valid {
			clipCount = 0
		}

		videos = append(videos, sbvision.Video{
			ID:        videoID,
			Title:     title,
			Thumbnail: sbvision.Image(s3Key),
			Duration:  videoDuration,
			FPS:       fps,
			ClipCount: clipCount,
		})
	}
	return videos, nil
}

// AddVideo adds the video to the database
func (db *SBDatabase) AddVideo(video *sbvision.Video) error {
	result, err := db.Exec(`
INSERT INTO videos (title, type, duration, fps, thumbnail_id) 
SELECT 
	?, ?, ?, ?, id
FROM images
WHERE images.key = ?
	`, video.Title, video.Type, video.Duration, video.FPS, video.Thumbnail)
	if err != nil {
		return err
	}
	video.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
