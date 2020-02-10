package database

import (
	"database/sql"
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddVideo() (err error) {
	sb.addVideo, err = sb.db.Prepare(`
INSERT INTO videos (title, type, format, duration, thumbnail_id) 
SELECT 
	?, ?, ?, ?, id
FROM images
WHERE images.key = ?
	`)
	return
}

// AddVideo adds the video to the database
func (sb *SBDatabase) AddVideo(video *sbvision.Video) error {
	result, err := sb.addVideo.Exec(video.Title, video.Type, video.Format, video.Duration, video.Thumbnail)
	if err != nil {
		return fmt.Errorf("\n\tError adding video: %s", err.Error())
	}
	video.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (sb *SBDatabase) prepareGetVideoCount() (err error) {
	sb.getVideoCount, err = sb.db.Prepare(`
SELECT COUNT(*) FROM videos
	`)
	return
}

// GetVideoCount gets the total number of videos in the database
func (sb *SBDatabase) GetVideoCount() (int64, error) {
	result := sb.getVideoCount.QueryRow()
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("\n\tError getting video count from db: %s", err.Error())
	}
	return count, nil
}

// prepareGetVideoById prepares the GetVideoById query
func (sb *SBDatabase) prepareGetVideoByID() (err error) {
	sb.getVideoByID, err = sb.db.Prepare(`
SELECT	
	videos.id,
	videos.title,
	images.key,
	videos.type,
	videos.format,
	videos.duration,
	COUNT(*),
	MAX(bounds.id)
FROM videos
INNER JOIN images 
		ON images.id = videos.thumbnail_id
LEFT JOIN frames
		ON frames.video_id = videos.id
LEFT JOIN bounds
		ON bounds.frame_id = frames.id
WHERE videos.id = ?
GROUP BY videos.id
ORDER BY discovery_time DESC
	`)
	return
}

// GetVideoByID gets a video record by it's id
func (sb *SBDatabase) GetVideoByID(id int64) (*sbvision.Video, error) {
	result := sb.getVideoByID.QueryRow(id)
	video := sbvision.Video{}
	err := sb.parseVideoRow(result, &video)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting video: %s", err)
	}
	return &video, nil
}

func (sb *SBDatabase) prepareGetVideos() (err error) {
	sb.getVideoPage, err = sb.db.Prepare(`
SELECT	
	videos.id,
	videos.title,
	images.key,
	videos.type,
	videos.format,
	videos.duration,
	COUNT(*),
	MAX(bounds.id)
FROM videos
INNER JOIN images 
		ON images.id = videos.thumbnail_id
LEFT JOIN frames
		ON frames.video_id = videos.id
LEFT JOIN bounds
		ON bounds.frame_id = frames.id
GROUP BY videos.id
ORDER BY discovery_time DESC
LIMIT ? OFFSET ?`)
	return
}

// GetVideos is a paged listing of videos
func (sb *SBDatabase) GetVideos(offset, count int64) ([]sbvision.Video, error) {
	results, err := sb.getVideoPage.Query(count, offset)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting database video records: %s", err.Error())
	}
	defer results.Close()
	var videos []sbvision.Video
	for results.Next() {
		videos = append(videos, sbvision.Video{})
		err = sb.parseVideoRow(results, &videos[len(videos)-1])
		if err != nil {
			return nil, fmt.Errorf("\n\tError Parsing video in list: %s", err.Error())
		}
	}
	return videos, nil
}

type scannable interface {
	Scan(to ...interface{}) error
}

func (sb *SBDatabase) parseVideoRow(src scannable, dst *sbvision.Video) error {
	var clipCount int64
	var clipFound sql.NullInt64

	err := src.Scan(
		&dst.ID,
		&dst.Title,
		&dst.Thumbnail,
		&dst.Type,
		&dst.Format,
		&dst.Duration,
		&clipCount,
		&clipFound,
	)
	if err != nil {
		return err
	}

	if !clipFound.Valid {
		dst.ClipCount = 0
	} else {
		dst.ClipCount = clipCount
	}

	return nil
}